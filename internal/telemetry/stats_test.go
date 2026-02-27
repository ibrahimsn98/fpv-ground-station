package telemetry

import (
	"strings"
	"testing"

	"fpv-ground-station/internal/ltm"
)

func TestStats_Count(t *testing.T) {
	s := NewStats()
	s.Count(ltm.FuncAttitude)
	s.Count(ltm.FuncAttitude)
	s.Count(ltm.FuncGPS)

	if s.Total != 3 {
		t.Errorf("total = %d, want 3", s.Total)
	}
	if s.Frames[ltm.FuncAttitude] != 2 {
		t.Errorf("attitude count = %d, want 2", s.Frames[ltm.FuncAttitude])
	}
	if s.Frames[ltm.FuncGPS] != 1 {
		t.Errorf("gps count = %d, want 1", s.Frames[ltm.FuncGPS])
	}
}

func TestStats_Errors(t *testing.T) {
	s := NewStats()
	s.RecordCRCError()
	s.RecordCRCError()
	s.RecordDecodeError()

	if s.CRCErrors != 2 {
		t.Errorf("crc errors = %d, want 2", s.CRCErrors)
	}
	if s.DecodeErrors != 1 {
		t.Errorf("decode errors = %d, want 1", s.DecodeErrors)
	}
}

func TestStats_FPS(t *testing.T) {
	s := NewStats()
	fps := s.FPS()
	if fps != 0 {
		t.Errorf("fps = %f, want 0", fps)
	}
}

func TestStats_Summary(t *testing.T) {
	s := NewStats()
	s.Count(ltm.FuncAttitude)
	s.Count(ltm.FuncGPS)
	s.RecordCRCError()

	summary := s.Summary()
	if !strings.Contains(summary, "LTM Telemetry Stats") {
		t.Error("summary missing header")
	}
	if !strings.Contains(summary, "Total:         2") {
		t.Error("summary missing total")
	}
	if !strings.Contains(summary, "CRC Errors:    1") {
		t.Error("summary missing CRC errors")
	}
	if !strings.Contains(summary, "Attitude") {
		t.Error("summary missing attitude frame")
	}
}

func TestStats_Snapshot(t *testing.T) {
	s := NewStats()
	s.Count(ltm.FuncAttitude)
	s.RecordCRCError()
	s.RecordDecodeError()

	snap := s.Snapshot()
	if snap.Total != 1 {
		t.Errorf("total = %d, want 1", snap.Total)
	}
	if snap.CRCErrors != 1 {
		t.Errorf("crc errors = %d, want 1", snap.CRCErrors)
	}
	if snap.DecodeErrors != 1 {
		t.Errorf("decode errors = %d, want 1", snap.DecodeErrors)
	}
}
