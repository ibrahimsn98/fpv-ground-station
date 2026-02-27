package ltm

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Decode interprets a RawFrame's payload and returns a typed Frame.
func Decode(f RawFrame) (Frame, error) {
	r := Frame{
		Function: f.Function,
		Name:     FrameName[f.Function],
		Time:     time.Now(),
	}

	var err error
	switch f.Function {
	case FuncGPS:
		r.GPS, err = decodeGPS(f.Payload)
	case FuncAttitude:
		r.Attitude, err = decodeAttitude(f.Payload)
	case FuncStatus:
		r.Status, err = decodeStatus(f.Payload)
	case FuncOrigin:
		r.Origin, err = decodeOrigin(f.Payload)
	case FuncNav:
		r.Nav, err = decodeNav(f.Payload)
	case FuncExtra:
		r.Extra, err = decodeExtra(f.Payload)
	default:
		return r, fmt.Errorf("ltm: unknown function 0x%02X", f.Function)
	}
	return r, err
}

func decodeGPS(p []byte) (*GPSData, error) {
	if len(p) < 14 {
		return nil, fmt.Errorf("ltm: gps payload too short (%d < 14)", len(p))
	}
	satInfo := p[13]
	return &GPSData{
		Lat:         float64(int32(binary.LittleEndian.Uint32(p[0:]))) / 1e7,
		Lon:         float64(int32(binary.LittleEndian.Uint32(p[4:]))) / 1e7,
		GroundSpeed: p[8],
		Altitude:    float64(int32(binary.LittleEndian.Uint32(p[9:]))) / 100.0,
		Fix:         satInfo & 0x03,
		Sats:        satInfo >> 2,
	}, nil
}

func decodeAttitude(p []byte) (*AttitudeData, error) {
	if len(p) < 6 {
		return nil, fmt.Errorf("ltm: attitude payload too short (%d < 6)", len(p))
	}
	return &AttitudeData{
		Pitch:   int16(binary.LittleEndian.Uint16(p[0:])),
		Roll:    int16(binary.LittleEndian.Uint16(p[2:])),
		Heading: int16(binary.LittleEndian.Uint16(p[4:])),
	}, nil
}

func decodeStatus(p []byte) (*StatusData, error) {
	if len(p) < 7 {
		return nil, fmt.Errorf("ltm: status payload too short (%d < 7)", len(p))
	}
	statusByte := p[6]
	return &StatusData{
		Vbat:       float64(binary.LittleEndian.Uint16(p[0:])) / 1000.0,
		MAhDrawn:   binary.LittleEndian.Uint16(p[2:]),
		RSSI:       p[4],
		Airspeed:   p[5],
		Armed:      statusByte&0x01 != 0,
		Failsafe:   statusByte&0x02 != 0,
		FlightMode: statusByte >> 2,
	}, nil
}

func decodeOrigin(p []byte) (*OriginData, error) {
	if len(p) < 14 {
		return nil, fmt.Errorf("ltm: origin payload too short (%d < 14)", len(p))
	}
	return &OriginData{
		Lat:   float64(int32(binary.LittleEndian.Uint32(p[0:]))) / 1e7,
		Lon:   float64(int32(binary.LittleEndian.Uint32(p[4:]))) / 1e7,
		Alt:   float64(int32(binary.LittleEndian.Uint32(p[8:]))) / 100.0,
		OSDOn: p[12]&0x01 != 0,
		Fix:   p[13],
	}, nil
}

func decodeNav(p []byte) (*NavData, error) {
	if len(p) < 6 {
		return nil, fmt.Errorf("ltm: nav payload too short (%d < 6)", len(p))
	}
	return &NavData{
		GPSMode:     p[0],
		NavMode:     p[1],
		NavAction:   p[2],
		WaypointNum: p[3],
		NavError:    p[4],
		Flags:       p[5],
	}, nil
}

func decodeExtra(p []byte) (*ExtraData, error) {
	if len(p) < 6 {
		return nil, fmt.Errorf("ltm: extra payload too short (%d < 6)", len(p))
	}
	return &ExtraData{
		HDOP:         float64(binary.LittleEndian.Uint16(p[0:])) / 100.0,
		HWStatus:     p[2],
		XCounter:     p[3],
		DisarmReason: p[4],
	}, nil
}
