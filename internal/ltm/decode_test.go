package ltm

import (
	"encoding/binary"
	"math"
	"testing"
)

func i16u(v int16) uint16 { return uint16(v) }
func i32u(v int32) uint32 { return uint32(v) }

func TestDecodeGPS(t *testing.T) {
	p := make([]byte, 14)
	binary.LittleEndian.PutUint32(p[0:], i32u(515000000))  // lat = 51.5°
	binary.LittleEndian.PutUint32(p[4:], i32u(-1278000))   // lon = -0.1278°
	p[8] = 15                                                // ground_speed = 15 m/s
	binary.LittleEndian.PutUint32(p[9:], i32u(10000))      // altitude = 100.00m
	p[13] = (12 << 2) | 3                                   // 12 sats, fix=3

	d, err := decodeGPS(p)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.Lat-51.5) > 0.0001 {
		t.Errorf("lat = %f, want 51.5", d.Lat)
	}
	if math.Abs(d.Lon-(-0.1278)) > 0.0001 {
		t.Errorf("lon = %f, want -0.1278", d.Lon)
	}
	if d.GroundSpeed != 15 {
		t.Errorf("ground_speed = %d, want 15", d.GroundSpeed)
	}
	if math.Abs(d.Altitude-100.0) > 0.01 {
		t.Errorf("altitude = %f, want 100.0", d.Altitude)
	}
	if d.Fix != 3 {
		t.Errorf("fix = %d, want 3", d.Fix)
	}
	if d.Sats != 12 {
		t.Errorf("sats = %d, want 12", d.Sats)
	}
}

func TestDecodeGPS_NegativeAltitude(t *testing.T) {
	p := make([]byte, 14)
	binary.LittleEndian.PutUint32(p[9:], i32u(-500)) // -5.00m
	p[13] = (5 << 2) | 2                              // 5 sats, fix=2

	d, err := decodeGPS(p)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.Altitude-(-5.0)) > 0.01 {
		t.Errorf("altitude = %f, want -5.0", d.Altitude)
	}
	if d.Fix != 2 {
		t.Errorf("fix = %d, want 2", d.Fix)
	}
	if d.Sats != 5 {
		t.Errorf("sats = %d, want 5", d.Sats)
	}
}

func TestDecodeGPS_MaxSats(t *testing.T) {
	p := make([]byte, 14)
	p[13] = (63 << 2) | 3 // max 63 sats, fix=3

	d, err := decodeGPS(p)
	if err != nil {
		t.Fatal(err)
	}
	if d.Sats != 63 {
		t.Errorf("sats = %d, want 63", d.Sats)
	}
}

func TestDecodeGPS_TooShort(t *testing.T) {
	_, err := decodeGPS(make([]byte, 10))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecodeAttitude(t *testing.T) {
	p := make([]byte, 6)
	binary.LittleEndian.PutUint16(p[0:], i16u(-15))  // pitch = -15°
	binary.LittleEndian.PutUint16(p[2:], i16u(30))   // roll = 30°
	binary.LittleEndian.PutUint16(p[4:], i16u(270))  // heading = 270°

	d, err := decodeAttitude(p)
	if err != nil {
		t.Fatal(err)
	}
	if d.Pitch != -15 {
		t.Errorf("pitch = %d, want -15", d.Pitch)
	}
	if d.Roll != 30 {
		t.Errorf("roll = %d, want 30", d.Roll)
	}
	if d.Heading != 270 {
		t.Errorf("heading = %d, want 270", d.Heading)
	}
}

func TestDecodeAttitude_TooShort(t *testing.T) {
	_, err := decodeAttitude(make([]byte, 3))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecodeStatus(t *testing.T) {
	p := make([]byte, 7)
	binary.LittleEndian.PutUint16(p[0:], 11800) // vbat = 11.8V
	binary.LittleEndian.PutUint16(p[2:], 1200)  // mAh drawn
	p[4] = 200                                   // rssi
	p[5] = 25                                    // airspeed
	p[6] = (10 << 2) | 0x03                     // flight mode 10 (Waypoints), armed + failsafe

	d, err := decodeStatus(p)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.Vbat-11.8) > 0.01 {
		t.Errorf("vbat = %f, want 11.8", d.Vbat)
	}
	if d.MAhDrawn != 1200 {
		t.Errorf("mah_drawn = %d, want 1200", d.MAhDrawn)
	}
	if d.RSSI != 200 {
		t.Errorf("rssi = %d, want 200", d.RSSI)
	}
	if d.Airspeed != 25 {
		t.Errorf("airspeed = %d, want 25", d.Airspeed)
	}
	if !d.Armed {
		t.Error("armed = false, want true")
	}
	if !d.Failsafe {
		t.Error("failsafe = false, want true")
	}
	if d.FlightMode != 10 {
		t.Errorf("flight_mode = %d, want 10", d.FlightMode)
	}
}

func TestDecodeStatus_Disarmed(t *testing.T) {
	p := make([]byte, 7)
	binary.LittleEndian.PutUint16(p[0:], 16200) // 16.2V
	p[6] = 2 << 2                                // mode 2 (Angle), not armed, not failsafe

	d, err := decodeStatus(p)
	if err != nil {
		t.Fatal(err)
	}
	if d.Armed {
		t.Error("armed = true, want false")
	}
	if d.Failsafe {
		t.Error("failsafe = true, want false")
	}
	if d.FlightMode != 2 {
		t.Errorf("flight_mode = %d, want 2", d.FlightMode)
	}
}

func TestDecodeStatus_AllFlightModes(t *testing.T) {
	for mode := uint8(0); mode <= 21; mode++ {
		p := make([]byte, 7)
		p[6] = mode << 2

		d, err := decodeStatus(p)
		if err != nil {
			t.Fatalf("mode %d: %v", mode, err)
		}
		if d.FlightMode != mode {
			t.Errorf("mode %d: flight_mode = %d", mode, d.FlightMode)
		}
	}
}

func TestDecodeStatus_TooShort(t *testing.T) {
	_, err := decodeStatus(make([]byte, 5))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecodeOrigin(t *testing.T) {
	p := make([]byte, 14)
	binary.LittleEndian.PutUint32(p[0:], i32u(515000000))  // home lat
	binary.LittleEndian.PutUint32(p[4:], i32u(-1278000))   // home lon
	binary.LittleEndian.PutUint32(p[8:], i32u(5000))       // home alt 50.00m
	p[12] = 0x01                                                     // osd_on
	p[13] = 3                                                        // fix

	d, err := decodeOrigin(p)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.Lat-51.5) > 0.0001 {
		t.Errorf("lat = %f, want 51.5", d.Lat)
	}
	if math.Abs(d.Lon-(-0.1278)) > 0.0001 {
		t.Errorf("lon = %f, want -0.1278", d.Lon)
	}
	if math.Abs(d.Alt-50.0) > 0.01 {
		t.Errorf("alt = %f, want 50.0", d.Alt)
	}
	if !d.OSDOn {
		t.Error("osd_on = false, want true")
	}
	if d.Fix != 3 {
		t.Errorf("fix = %d, want 3", d.Fix)
	}
}

func TestDecodeOrigin_TooShort(t *testing.T) {
	_, err := decodeOrigin(make([]byte, 10))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecodeNav(t *testing.T) {
	p := []byte{2, 5, 1, 7, 0, 0xFF}

	d, err := decodeNav(p)
	if err != nil {
		t.Fatal(err)
	}
	if d.GPSMode != 2 {
		t.Errorf("gps_mode = %d, want 2", d.GPSMode)
	}
	if d.NavMode != 5 {
		t.Errorf("nav_mode = %d, want 5", d.NavMode)
	}
	if d.NavAction != 1 {
		t.Errorf("nav_action = %d, want 1", d.NavAction)
	}
	if d.WaypointNum != 7 {
		t.Errorf("waypoint_num = %d, want 7", d.WaypointNum)
	}
	if d.NavError != 0 {
		t.Errorf("nav_error = %d, want 0", d.NavError)
	}
	if d.Flags != 0xFF {
		t.Errorf("flags = 0x%02X, want 0xFF", d.Flags)
	}
}

func TestDecodeNav_TooShort(t *testing.T) {
	_, err := decodeNav(make([]byte, 3))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecodeExtra(t *testing.T) {
	p := make([]byte, 6)
	binary.LittleEndian.PutUint16(p[0:], 150) // hdop = 1.50
	p[2] = 0                                   // hw_status OK
	p[3] = 42                                  // x_counter
	p[4] = 3                                   // disarm_reason

	d, err := decodeExtra(p)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.HDOP-1.50) > 0.01 {
		t.Errorf("hdop = %f, want 1.50", d.HDOP)
	}
	if d.HWStatus != 0 {
		t.Errorf("hw_status = %d, want 0", d.HWStatus)
	}
	if d.XCounter != 42 {
		t.Errorf("x_counter = %d, want 42", d.XCounter)
	}
	if d.DisarmReason != 3 {
		t.Errorf("disarm_reason = %d, want 3", d.DisarmReason)
	}
}

func TestDecodeExtra_TooShort(t *testing.T) {
	_, err := decodeExtra(make([]byte, 3))
	if err == nil {
		t.Error("expected error for short payload")
	}
}

func TestDecode_Dispatch(t *testing.T) {
	tests := []struct {
		fn   byte
		size int
	}{
		{FuncGPS, 14},
		{FuncAttitude, 6},
		{FuncStatus, 7},
		{FuncOrigin, 14},
		{FuncNav, 6},
		{FuncExtra, 6},
	}
	for _, tt := range tests {
		name := FrameName[tt.fn]
		t.Run(name, func(t *testing.T) {
			f := RawFrame{Function: tt.fn, Payload: make([]byte, tt.size)}
			r, err := Decode(f)
			if err != nil {
				t.Fatalf("Decode(%s): %v", name, err)
			}
			if r.Function != tt.fn {
				t.Errorf("function = 0x%02X, want 0x%02X", r.Function, tt.fn)
			}
			if r.Name != name {
				t.Errorf("name = %q, want %q", r.Name, name)
			}
		})
	}
}

func TestDecode_UnknownFunction(t *testing.T) {
	f := RawFrame{Function: 'Z', Payload: nil}
	_, err := Decode(f)
	if err == nil {
		t.Error("expected error for unknown function")
	}
}
