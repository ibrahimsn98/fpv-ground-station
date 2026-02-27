package telemetry

import (
	"sync"
	"time"

	"fpv-ground-station/internal/ltm"
)

// Store holds the latest telemetry state, safe for concurrent access.
type Store struct {
	mu sync.RWMutex

	GPS          *ltm.GPSData
	GPSTime      time.Time
	Attitude     *ltm.AttitudeData
	AttitudeTime time.Time
	Status       *ltm.StatusData
	StatusTime   time.Time
	Origin       *ltm.OriginData
	OriginTime   time.Time
	Nav          *ltm.NavData
	NavTime      time.Time
	Extra        *ltm.ExtraData
	ExtraTime    time.Time
}

// Snapshot is a point-in-time copy of telemetry state, safe to use without locks.
type Snapshot struct {
	GPS          *ltm.GPSData
	GPSTime      time.Time
	Attitude     *ltm.AttitudeData
	AttitudeTime time.Time
	Status       *ltm.StatusData
	StatusTime   time.Time
	Origin       *ltm.OriginData
	OriginTime   time.Time
	Nav          *ltm.NavData
	NavTime      time.Time
	Extra        *ltm.ExtraData
	ExtraTime    time.Time
}

// Update merges a decoded LTM frame into the store.
func (s *Store) Update(f ltm.Frame) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch {
	case f.GPS != nil:
		s.GPS = f.GPS
		s.GPSTime = f.Time
	case f.Attitude != nil:
		s.Attitude = f.Attitude
		s.AttitudeTime = f.Time
	case f.Status != nil:
		s.Status = f.Status
		s.StatusTime = f.Time
	case f.Origin != nil:
		s.Origin = f.Origin
		s.OriginTime = f.Time
	case f.Nav != nil:
		s.Nav = f.Nav
		s.NavTime = f.Time
	case f.Extra != nil:
		s.Extra = f.Extra
		s.ExtraTime = f.Time
	}
}

// Snapshot returns a point-in-time copy of the current telemetry state.
func (s *Store) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return Snapshot{
		GPS:          s.GPS,
		GPSTime:      s.GPSTime,
		Attitude:     s.Attitude,
		AttitudeTime: s.AttitudeTime,
		Status:       s.Status,
		StatusTime:   s.StatusTime,
		Origin:       s.Origin,
		OriginTime:   s.OriginTime,
		Nav:          s.Nav,
		NavTime:      s.NavTime,
		Extra:        s.Extra,
		ExtraTime:    s.ExtraTime,
	}
}
