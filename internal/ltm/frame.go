package ltm

import "time"

// LTM frame header bytes.
const (
	Header1 = '$'
	Header2 = 'T'
)

// Frame function characters.
const (
	FuncGPS      byte = 'G'
	FuncAttitude byte = 'A'
	FuncStatus   byte = 'S'
	FuncOrigin   byte = 'O'
	FuncNav      byte = 'N'
	FuncExtra    byte = 'X'
)

// PayloadSize maps each frame function to its expected payload size.
var PayloadSize = map[byte]int{
	FuncGPS:      14,
	FuncAttitude: 6,
	FuncStatus:   7,
	FuncOrigin:   14,
	FuncNav:      6,
	FuncExtra:    6,
}

// FrameName maps each frame function to a human-readable name.
var FrameName = map[byte]string{
	FuncGPS:      "GPS",
	FuncAttitude: "Attitude",
	FuncStatus:   "Status",
	FuncOrigin:   "Origin",
	FuncNav:      "Nav",
	FuncExtra:    "Extra",
}

// Flight mode enum (0-21).
var FlightModeName = map[uint8]string{
	0:  "Manual",
	1:  "Rate",
	2:  "Angle",
	3:  "Horizon",
	4:  "Acro",
	5:  "Stabilized1",
	6:  "Stabilized2",
	7:  "Stabilized3",
	8:  "Altitude Hold",
	9:  "GPS Hold",
	10: "Waypoints",
	11: "Head Free",
	12: "Circle",
	13: "RTH",
	14: "Follow Me",
	15: "Land",
	16: "Fly By Wire A",
	17: "Fly By Wire B",
	18: "Cruise",
	19: "Unknown",
	20: "Launch",
	21: "Autotune",
}

// GPS mode enum (0-3).
var GPSModeName = map[uint8]string{
	0: "None",
	1: "PosHold",
	2: "RTH",
	3: "Mission",
}

// Nav mode enum (0-15).
var NavModeName = map[uint8]string{
	0:  "None",
	1:  "RTH Start",
	2:  "RTH Enroute",
	3:  "PosHold Infinite",
	4:  "PosHold Timed",
	5:  "WP Enroute",
	6:  "Process Next",
	7:  "Jump",
	8:  "Start Land",
	9:  "Land In Progress",
	10: "Landed",
	11: "Settling Before Land",
	12: "Start Descent",
	13: "Hover Above Home",
	14: "Emergency Landing",
	15: "Critical GPS",
}

// Nav action enum (0-8).
var NavActionName = map[uint8]string{
	0: "Unassigned",
	1: "Waypoint",
	2: "PosHold Unlim",
	3: "PosHold Timed",
	4: "RTH",
	5: "Set POI",
	6: "Jump",
	7: "Set Head",
	8: "Land",
}

// Nav error enum (0-11).
var NavErrorName = map[uint8]string{
	0:  "OK",
	1:  "WP Too Far",
	2:  "WP Sanity",
	3:  "WP CRC",
	4:  "Finish",
	5:  "Timer Complete",
	6:  "Invalid Jump",
	7:  "Invalid Data",
	8:  "Wait For RTH Alt",
	9:  "GPS Fix Lost",
	10: "Disarmed",
	11: "Landing Check",
}

// GPSData from the G-frame (14 bytes, 5 Hz).
type GPSData struct {
	Lat         float64 `json:"lat"`          // degrees (int32/1e7)
	Lon         float64 `json:"lon"`          // degrees (int32/1e7)
	GroundSpeed uint8   `json:"ground_speed"` // m/s
	Altitude    float64 `json:"altitude"`     // meters (int32 cm / 100)
	Fix         uint8   `json:"fix"`          // 0=no fix, 1=dead reck, 2=2D, 3=3D
	Sats        uint8   `json:"sats"`         // satellite count
}

// AttitudeData from the A-frame (6 bytes, 10 Hz).
type AttitudeData struct {
	Pitch   int16 `json:"pitch"`   // degrees
	Roll    int16 `json:"roll"`    // degrees
	Heading int16 `json:"heading"` // degrees
}

// StatusData from the S-frame (7 bytes, 5 Hz).
type StatusData struct {
	Vbat       float64 `json:"vbat"`        // volts (uint16 mV / 1000)
	MAhDrawn   uint16  `json:"mah_drawn"`   // consumed mAh
	RSSI       uint8   `json:"rssi"`        // 0-254
	Airspeed   uint8   `json:"airspeed"`    // m/s
	Armed      bool    `json:"armed"`       // bit 0 of status byte
	Failsafe   bool    `json:"failsafe"`    // bit 1 of status byte
	FlightMode uint8   `json:"flight_mode"` // enum 0-21 (bits 2-7 >> 2)
}

// OriginData from the O-frame (14 bytes, 1 Hz).
type OriginData struct {
	Lat   float64 `json:"lat"`    // degrees (int32/1e7)
	Lon   float64 `json:"lon"`    // degrees (int32/1e7)
	Alt   float64 `json:"alt"`    // meters (int32 cm / 100)
	OSDOn bool    `json:"osd_on"` // byte[12] bit 0
	Fix   uint8   `json:"fix"`    // byte[13]
}

// NavData from the N-frame (6 bytes, ~4 Hz).
type NavData struct {
	GPSMode    uint8 `json:"gps_mode"`    // enum 0-3
	NavMode    uint8 `json:"nav_mode"`    // enum 0-15
	NavAction  uint8 `json:"nav_action"`  // enum 0-8
	WaypointNum uint8 `json:"waypoint_num"`
	NavError   uint8 `json:"nav_error"`   // enum 0-11
	Flags      uint8 `json:"flags"`
}

// ExtraData from the X-frame (6 bytes, 1 Hz).
type ExtraData struct {
	HDOP          float64 `json:"hdop"`           // (uint16 / 100)
	HWStatus      uint8   `json:"hw_status"`      // 0=OK, 1=fail
	XCounter      uint8   `json:"x_counter"`      // packet counter
	DisarmReason  uint8   `json:"disarm_reason"`
}

// RawFrame is a validated LTM frame before payload decoding.
type RawFrame struct {
	Function byte
	Payload  []byte
}

// Frame is a decoded LTM frame. Exactly one data pointer is non-nil.
type Frame struct {
	Function byte      `json:"function"`
	Name     string    `json:"name"`
	Time     time.Time `json:"time"`

	GPS      *GPSData      `json:"gps,omitempty"`
	Attitude *AttitudeData `json:"attitude,omitempty"`
	Status   *StatusData   `json:"status,omitempty"`
	Origin   *OriginData   `json:"origin,omitempty"`
	Nav      *NavData      `json:"nav,omitempty"`
	Extra    *ExtraData    `json:"extra,omitempty"`
}
