package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"fpv-ground-station/internal/ltm"
	"fpv-ground-station/internal/telemetry"

	"nhooyr.io/websocket"
)

// Message is the JSON envelope sent to each WebSocket client.
type Message struct {
	Timestamp int64 `json:"ts"` // Unix millis

	GPS          *ltm.GPSData      `json:"gps,omitempty"`
	GPSTime      int64             `json:"gps_ts,omitempty"`
	Attitude     *ltm.AttitudeData `json:"attitude,omitempty"`
	AttitudeTime int64             `json:"attitude_ts,omitempty"`
	Status       *ltm.StatusData   `json:"status,omitempty"`
	StatusTime   int64             `json:"status_ts,omitempty"`
	Origin       *ltm.OriginData   `json:"origin,omitempty"`
	OriginTime   int64             `json:"origin_ts,omitempty"`
	Nav          *ltm.NavData      `json:"nav,omitempty"`
	NavTime      int64             `json:"nav_ts,omitempty"`
	Extra        *ltm.ExtraData    `json:"extra,omitempty"`
	ExtraTime    int64             `json:"extra_ts,omitempty"`

	Stats *StatsPayload `json:"stats,omitempty"`
}

// StatsPayload contains connection/throughput metrics.
type StatsPayload struct {
	UptimeSec    float64 `json:"uptime_sec"`
	Total        int     `json:"total"`
	FPS          float64 `json:"fps"`
	CRCErrors    int     `json:"crc_errors"`
	DecodeErrors int     `json:"decode_errors"`
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow any origin (Vite dev server)
	})
	if err != nil {
		log.Printf("ws accept: %v", err)
		return
	}

	c := &client{send: make(chan []byte, 16)}
	s.addClient(c)

	ctx := r.Context()

	// Writer goroutine
	go func() {
		defer conn.Close(websocket.StatusNormalClosure, "")
		defer s.removeClient(c)

		for {
			select {
			case msg, ok := <-c.send:
				if !ok {
					return
				}
				err := conn.Write(ctx, websocket.MessageText, msg)
				if err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Reader goroutine — just drain incoming messages to detect close
	for {
		_, _, err := conn.Read(ctx)
		if err != nil {
			break
		}
	}
}

func (s *Server) broadcastLoop(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := s.buildMessage()
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("ws marshal: %v", err)
				continue
			}

			s.mu.RLock()
			for c := range s.clients {
				select {
				case c.send <- data:
				default:
					// slow client, drop frame
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *Server) buildMessage() Message {
	snap := s.store.Snapshot()
	statsSnap := s.stats.Snapshot()

	msg := Message{
		Timestamp: time.Now().UnixMilli(),
		Stats: &StatsPayload{
			UptimeSec:    statsSnap.UptimeSec,
			Total:        statsSnap.Total,
			FPS:          statsSnap.FPS,
			CRCErrors:    statsSnap.CRCErrors,
			DecodeErrors: statsSnap.DecodeErrors,
		},
	}

	if snap.GPS != nil {
		msg.GPS = snap.GPS
		msg.GPSTime = toMillis(snap.GPSTime)
	}
	if snap.Attitude != nil {
		msg.Attitude = snap.Attitude
		msg.AttitudeTime = toMillis(snap.AttitudeTime)
	}
	if snap.Status != nil {
		msg.Status = snap.Status
		msg.StatusTime = toMillis(snap.StatusTime)
	}
	if snap.Origin != nil {
		msg.Origin = snap.Origin
		msg.OriginTime = toMillis(snap.OriginTime)
	}
	if snap.Nav != nil {
		msg.Nav = snap.Nav
		msg.NavTime = toMillis(snap.NavTime)
	}
	if snap.Extra != nil {
		msg.Extra = snap.Extra
		msg.ExtraTime = toMillis(snap.ExtraTime)
	}

	return msg
}

func toMillis(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

// statsFromTelemetry converts a telemetry StatsSnapshot to a StatsPayload.
func statsFromTelemetry(snap telemetry.StatsSnapshot) StatsPayload {
	return StatsPayload{
		UptimeSec:    snap.UptimeSec,
		Total:        snap.Total,
		FPS:          snap.FPS,
		CRCErrors:    snap.CRCErrors,
		DecodeErrors: snap.DecodeErrors,
	}
}
