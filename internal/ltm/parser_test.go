package ltm

import (
	"bytes"
	"testing"
)

// buildLTMFrame constructs a valid LTM frame: $T + func + payload + checksum
func buildLTMFrame(fn byte, payload []byte) []byte {
	frame := make([]byte, 0, 3+len(payload)+1)
	frame = append(frame, Header1, Header2, fn)
	frame = append(frame, payload...)
	frame = append(frame, xorChecksum(payload))
	return frame
}

func TestXorChecksum_KnownVectors(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		want    byte
	}{
		{"empty", nil, 0x00},
		{"one_byte", []byte{0xFF}, 0xFF},
		{"two_bytes", []byte{0x01, 0x02}, 0x03},
		{"all_zeros", make([]byte, 14), 0x00},
		{"real_attitude", []byte{0x01, 0x00, 0x00, 0x00, 0x67, 0x01}, 0x67},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xorChecksum(tt.payload)
			if got != tt.want {
				t.Errorf("xorChecksum(%x) = 0x%02X, want 0x%02X", tt.payload, got, tt.want)
			}
		})
	}
}

func TestParser_SingleFrame(t *testing.T) {
	payload := make([]byte, 6)
	payload[0] = 0x05 // pitch lo
	payload[1] = 0x00 // pitch hi
	payload[2] = 0x0A // roll lo
	payload[3] = 0x00 // roll hi
	payload[4] = 0xB4 // heading lo (180)
	payload[5] = 0x00 // heading hi
	frame := buildLTMFrame(FuncAttitude, payload)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(frame)

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
	if got[0].Function != FuncAttitude {
		t.Errorf("function = 0x%02X, want 0x%02X", got[0].Function, FuncAttitude)
	}
	if !bytes.Equal(got[0].Payload, payload) {
		t.Errorf("payload = %x, want %x", got[0].Payload, payload)
	}
}

func TestParser_ByteAtATime(t *testing.T) {
	payload := make([]byte, 6)
	frame := buildLTMFrame(FuncAttitude, payload)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	for _, b := range frame {
		p.Write([]byte{b})
	}

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
}

func TestParser_ConsecutiveFrames(t *testing.T) {
	att := buildLTMFrame(FuncAttitude, make([]byte, 6))
	gps := buildLTMFrame(FuncGPS, make([]byte, 14))

	var combined []byte
	combined = append(combined, att...)
	combined = append(combined, gps...)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(combined)

	if len(got) != 2 {
		t.Fatalf("got %d frames, want 2", len(got))
	}
	if got[0].Function != FuncAttitude {
		t.Errorf("frame[0] function = 0x%02X, want Attitude", got[0].Function)
	}
	if got[1].Function != FuncGPS {
		t.Errorf("frame[1] function = 0x%02X, want GPS", got[1].Function)
	}
}

func TestParser_GarbageTolerance(t *testing.T) {
	frame := buildLTMFrame(FuncStatus, make([]byte, 7))

	garbage := []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA}
	var data []byte
	data = append(data, garbage...)
	data = append(data, frame...)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(data)

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
	if got[0].Function != FuncStatus {
		t.Errorf("function = 0x%02X, want Status", got[0].Function)
	}
}

func TestParser_ChecksumFailureAndResync(t *testing.T) {
	bad := buildLTMFrame(FuncAttitude, make([]byte, 6))
	bad[len(bad)-1] ^= 0xFF // corrupt checksum

	good := buildLTMFrame(FuncGPS, make([]byte, 14))

	var data []byte
	data = append(data, bad...)
	data = append(data, good...)

	var frames []RawFrame
	var errs []error
	p := NewParser(func(f RawFrame) {
		frames = append(frames, f)
	}, func(e error) {
		errs = append(errs, e)
	})

	p.Write(data)

	if len(errs) != 1 {
		t.Errorf("got %d errors, want 1", len(errs))
	}
	if len(frames) != 1 {
		t.Fatalf("got %d frames, want 1", len(frames))
	}
	if frames[0].Function != FuncGPS {
		t.Errorf("function = 0x%02X, want GPS", frames[0].Function)
	}
}

func TestParser_UnknownFunction(t *testing.T) {
	// '$' 'T' 'Z' — unknown function
	data := []byte{Header1, Header2, 'Z'}
	good := buildLTMFrame(FuncAttitude, make([]byte, 6))
	data = append(data, good...)

	var frames []RawFrame
	var errs []error
	p := NewParser(func(f RawFrame) {
		frames = append(frames, f)
	}, func(e error) {
		errs = append(errs, e)
	})

	p.Write(data)

	if len(errs) != 1 {
		t.Errorf("got %d errors, want 1", len(errs))
	}
	if len(frames) != 1 {
		t.Fatalf("got %d frames, want 1", len(frames))
	}
}

func TestParser_PartialWrites(t *testing.T) {
	frame := buildLTMFrame(FuncGPS, make([]byte, 14))

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	// Feed in chunks of 3 bytes
	for i := 0; i < len(frame); i += 3 {
		end := i + 3
		if end > len(frame) {
			end = len(frame)
		}
		p.Write(frame[i:end])
	}

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
}

func TestParser_FalseDollarSign(t *testing.T) {
	// '$' followed by not-'T' should not break parser
	data := []byte{'$', 'X'}
	frame := buildLTMFrame(FuncAttitude, make([]byte, 6))
	data = append(data, frame...)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(data)

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
}

func TestParser_AllFrameTypes(t *testing.T) {
	types := []struct {
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

	var data []byte
	for _, tt := range types {
		data = append(data, buildLTMFrame(tt.fn, make([]byte, tt.size))...)
	}

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(data)

	if len(got) != len(types) {
		t.Fatalf("got %d frames, want %d", len(got), len(types))
	}
	for i, tt := range types {
		if got[i].Function != tt.fn {
			t.Errorf("frame[%d] function = 0x%02X, want 0x%02X", i, got[i].Function, tt.fn)
		}
	}
}

func TestParser_DollarInPayload(t *testing.T) {
	// '$' appearing in payload should not disrupt parsing
	payload := make([]byte, 6)
	payload[2] = '$' // '$' in middle of attitude payload
	frame := buildLTMFrame(FuncAttitude, payload)

	var got []RawFrame
	p := NewParser(func(f RawFrame) {
		got = append(got, f)
	}, nil)

	p.Write(frame)

	if len(got) != 1 {
		t.Fatalf("got %d frames, want 1", len(got))
	}
}
