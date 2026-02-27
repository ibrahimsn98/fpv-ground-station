package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fpv-ground-station/internal/ltm"
	"fpv-ground-station/internal/serial"
	"fpv-ground-station/internal/server"
	"fpv-ground-station/internal/telemetry"
)

func main() {
	portName := flag.String("port", envOr("PORT", "/dev/cu.usbserial-840"), "serial port path")
	flag.StringVar(portName, "p", *portName, "serial port path (shorthand)")
	baud := flag.Int("baud", envOrInt("BAUD", 19200), "baud rate")
	flag.IntVar(baud, "b", *baud, "baud rate (shorthand)")
	jsonOut := flag.Bool("json", false, "output JSON lines instead of human-readable")
	webAddr := flag.String("web", ":8080", "web UI listen address (e.g. :8080)")
	devMode := flag.Bool("dev", false, "dev mode: skip embedded UI, use Vite proxy")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port, err := serial.Open(serial.Config{Name: *portName, Baud: *baud})
	if err != nil {
		log.Fatalf("open serial: %v", err)
	}
	defer port.Close()

	// Flush stale bytes before starting
	port.ResetInputBuffer()

	log.Printf("LTM on %s @ %d baud", *portName, *baud)

	store := &telemetry.Store{}
	stats := telemetry.NewStats()

	// Start web server
	distFS, err := webDistFS()
	if err != nil {
		log.Fatalf("load embedded UI: %v", err)
	}
	if distFS == nil && !*devMode {
		log.Fatal("no embedded UI available; rebuild with 'make build' or use --dev flag")
	}

	srv := server.New(server.Config{
		Store:   store,
		Stats:   stats,
		Addr:    *webAddr,
		WebFS:   distFS,
		DevMode: *devMode,
	})

	go func() {
		if err := srv.ListenAndServe(ctx); err != nil {
			log.Fatalf("web server: %v", err)
		}
	}()

	log.Printf("Web UI: http://localhost%s", *webAddr)

	// Perf ticker: log attitude Hz every second
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				hz := stats.AttitudeRx.Swap(0)
				if hz > 0 {
					log.Printf("Attitude: %d Hz", hz)
				}
			}
		}
	}()

	// Read LTM from serial port
	readLTM(ctx, port, store, stats, *jsonOut)

	// Shutdown
	port.Close()
	log.Println("Shutting down...")

	if *jsonOut {
		statsJSON := struct {
			UptimeSec    float64      `json:"uptime_sec"`
			Total        int          `json:"total"`
			FPS          float64      `json:"fps"`
			Frames       map[byte]int `json:"frames"`
			CRCErrors    int          `json:"crc_errors"`
			DecodeErrors int          `json:"decode_errors"`
		}{
			UptimeSec:    stats.Uptime().Seconds(),
			Total:        stats.Total,
			FPS:          stats.FPS(),
			Frames:       stats.Frames,
			CRCErrors:    stats.CRCErrors,
			DecodeErrors: stats.DecodeErrors,
		}
		json.NewEncoder(os.Stderr).Encode(statsJSON)
	} else {
		fmt.Fprintln(os.Stderr, stats.Summary())
	}
}

func readLTM(ctx context.Context, port *serial.Port, store *telemetry.Store, stats *telemetry.Stats, jsonOut bool) {
	enc := json.NewEncoder(os.Stdout)

	parser := ltm.NewParser(
		func(raw ltm.RawFrame) {
			frame, err := ltm.Decode(raw)
			if err != nil {
				stats.RecordDecodeError()
				return
			}

			store.Update(frame)
			stats.Count(frame.Function)

			if frame.Attitude != nil {
				stats.AttitudeRx.Add(1)
			}

			if jsonOut {
				enc.Encode(frame)
			} else {
				printHuman(frame)
			}
		},
		func(err error) {
			log.Printf("[PARSER ERR] %v", err)
			stats.RecordCRCError()
		},
	)

	buf := make([]byte, 256)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := port.Read(buf)
		if err != nil {
			continue
		}
		if n > 0 {
			parser.Write(buf[:n])
		}
	}
}

func printHuman(f ltm.Frame) {
	switch {
	case f.Attitude != nil:
		a := f.Attitude
		fmt.Printf("[ATT] pitch=%d roll=%d heading=%d\n",
			a.Pitch, a.Roll, a.Heading)
	case f.GPS != nil:
		g := f.GPS
		fmt.Printf("[GPS] lat=%.7f lon=%.7f alt=%.1fm spd=%dm/s fix=%d sats=%d\n",
			g.Lat, g.Lon, g.Altitude, g.GroundSpeed, g.Fix, g.Sats)
	case f.Status != nil:
		s := f.Status
		armed := "DISARMED"
		if s.Armed {
			armed = "ARMED"
		}
		mode := ltm.FlightModeName[s.FlightMode]
		fmt.Printf("[STS] %s mode=%s vbat=%.2fV mAh=%d rssi=%d airspeed=%d failsafe=%v\n",
			armed, mode, s.Vbat, s.MAhDrawn, s.RSSI, s.Airspeed, s.Failsafe)
	case f.Origin != nil:
		o := f.Origin
		fmt.Printf("[ORI] home_lat=%.7f home_lon=%.7f home_alt=%.1fm fix=%d\n",
			o.Lat, o.Lon, o.Alt, o.Fix)
	case f.Nav != nil:
		n := f.Nav
		fmt.Printf("[NAV] gps_mode=%d nav_mode=%d action=%d wp=%d error=%d\n",
			n.GPSMode, n.NavMode, n.NavAction, n.WaypointNum, n.NavError)
	case f.Extra != nil:
		x := f.Extra
		fmt.Printf("[EXT] hdop=%.2f hw=%d counter=%d disarm_reason=%d\n",
			x.HDOP, x.HWStatus, x.XCounter, x.DisarmReason)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envOrInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		var x int
		_, _ = fmt.Sscanf(v, "%d", &x)
		if x > 0 {
			return x
		}
	}
	return def
}
