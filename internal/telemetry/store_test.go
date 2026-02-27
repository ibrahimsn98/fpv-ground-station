package telemetry

import (
	"sync"
	"testing"
	"time"

	"fpv-ground-station/internal/ltm"
)

func TestStore_Update_Attitude(t *testing.T) {
	s := &Store{}
	f := ltm.Frame{
		Function: ltm.FuncAttitude,
		Time:     time.Now(),
		Attitude: &ltm.AttitudeData{
			Roll: 10, Pitch: -5, Heading: 180,
		},
	}

	s.Update(f)
	snap := s.Snapshot()

	if snap.Attitude == nil {
		t.Fatal("expected non-nil Attitude")
	}
	if snap.Attitude.Roll != 10 {
		t.Errorf("roll = %d, want 10", snap.Attitude.Roll)
	}
}

func TestStore_Update_AllTypes(t *testing.T) {
	s := &Store{}
	now := time.Now()

	updates := []ltm.Frame{
		{Function: ltm.FuncGPS, Time: now, GPS: &ltm.GPSData{Fix: 3}},
		{Function: ltm.FuncAttitude, Time: now, Attitude: &ltm.AttitudeData{Roll: 1}},
		{Function: ltm.FuncStatus, Time: now, Status: &ltm.StatusData{Vbat: 11.8}},
		{Function: ltm.FuncOrigin, Time: now, Origin: &ltm.OriginData{Lat: 51.5}},
		{Function: ltm.FuncNav, Time: now, Nav: &ltm.NavData{GPSMode: 2}},
		{Function: ltm.FuncExtra, Time: now, Extra: &ltm.ExtraData{HDOP: 1.5}},
	}

	for _, f := range updates {
		s.Update(f)
	}

	snap := s.Snapshot()
	if snap.GPS == nil || snap.Attitude == nil || snap.Status == nil ||
		snap.Origin == nil || snap.Nav == nil || snap.Extra == nil {
		t.Error("expected all fields to be non-nil after full update")
	}
}

func TestStore_Overwrite(t *testing.T) {
	s := &Store{}
	s.Update(ltm.Frame{
		Function: ltm.FuncAttitude,
		Time:     time.Now(),
		Attitude: &ltm.AttitudeData{Roll: 5},
	})
	s.Update(ltm.Frame{
		Function: ltm.FuncAttitude,
		Time:     time.Now(),
		Attitude: &ltm.AttitudeData{Roll: 15},
	})

	snap := s.Snapshot()
	if snap.Attitude.Roll != 15 {
		t.Errorf("roll = %d, want 15 (latest)", snap.Attitude.Roll)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s := &Store{}
	var wg sync.WaitGroup

	// 50 concurrent writers
	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.Update(ltm.Frame{
				Function: ltm.FuncAttitude,
				Time:     time.Now(),
				Attitude: &ltm.AttitudeData{Roll: int16(n)},
			})
		}(i)
	}

	// 50 concurrent readers
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Snapshot()
		}()
	}

	wg.Wait()
	// No race condition — test passes if no panic
}
