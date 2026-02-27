package telemetry

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fpv-ground-station/internal/ltm"
)

// Stats tracks telemetry metrics: frame counts, errors, and throughput.
type Stats struct {
	mu sync.Mutex

	Frames       map[byte]int
	CRCErrors    int
	DecodeErrors int
	Total        int
	StartTime    time.Time

	// Attitude receive rate counter (reset every second by perf ticker)
	AttitudeRx atomic.Int64
}

// NewStats creates a Stats tracker starting now.
func NewStats() *Stats {
	return &Stats{
		Frames:    make(map[byte]int),
		StartTime: time.Now(),
	}
}

// Count records a successfully decoded frame.
func (s *Stats) Count(fn byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Frames[fn]++
	s.Total++
}

// RecordCRCError increments the CRC error counter.
func (s *Stats) RecordCRCError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CRCErrors++
}

// RecordDecodeError increments the decode error counter.
func (s *Stats) RecordDecodeError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DecodeErrors++
}

// Uptime returns the duration since tracking started.
func (s *Stats) Uptime() time.Duration {
	return time.Since(s.StartTime)
}

// FPS returns the average frames per second.
func (s *Stats) FPS() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	sec := time.Since(s.StartTime).Seconds()
	if sec <= 0 {
		return 0
	}
	return float64(s.Total) / sec
}

// StatsSnapshot is a point-in-time copy of stats, safe to use without locks.
type StatsSnapshot struct {
	UptimeSec    float64
	Total        int
	FPS          float64
	CRCErrors    int
	DecodeErrors int
}

// Snapshot returns a thread-safe copy of all stat counters.
func (s *Stats) Snapshot() StatsSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	sec := time.Since(s.StartTime).Seconds()
	fps := 0.0
	if sec > 0 {
		fps = float64(s.Total) / sec
	}
	return StatsSnapshot{
		UptimeSec:    sec,
		Total:        s.Total,
		FPS:          fps,
		CRCErrors:    s.CRCErrors,
		DecodeErrors: s.DecodeErrors,
	}
}

// Summary returns a human-readable summary of all metrics.
func (s *Stats) Summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var b strings.Builder
	sec := time.Since(s.StartTime).Seconds()
	fps := 0.0
	if sec > 0 {
		fps = float64(s.Total) / sec
	}

	fmt.Fprintf(&b, "--- LTM Telemetry Stats ---\n")
	fmt.Fprintf(&b, "Uptime:        %.1fs\n", sec)
	fmt.Fprintf(&b, "Total:         %d\n", s.Total)
	fmt.Fprintf(&b, "FPS:           %.1f\n", fps)
	fmt.Fprintf(&b, "CRC Errors:    %d\n", s.CRCErrors)
	fmt.Fprintf(&b, "Decode Errors: %d\n", s.DecodeErrors)

	if len(s.Frames) > 0 {
		fmt.Fprintf(&b, "Frames:\n")
		for fn, count := range s.Frames {
			name := ltm.FrameName[fn]
			if name == "" {
				name = fmt.Sprintf("0x%02X", fn)
			}
			fmt.Fprintf(&b, "  %-25s %d\n", name, count)
		}
	}
	return b.String()
}
