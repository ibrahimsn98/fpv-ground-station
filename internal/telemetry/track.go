package telemetry

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// TrackLog persists GPS coordinates to a CSV file (lat,lon per line).
type TrackLog struct {
	mu   sync.Mutex
	path string
	file *os.File
}

// NewTrackLog opens or creates the track file in append mode.
func NewTrackLog(path string) (*TrackLog, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &TrackLog{path: path, file: f}, nil
}

// Append writes a single coordinate to the track file.
func (t *TrackLog) Append(lat, lon float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, err := fmt.Fprintf(t.file, "%.7f,%.7f\n", lat, lon)
	return err
}

// ReadAll reads all coordinates from the track file.
func (t *TrackLog) ReadAll() ([][2]float64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var points [][2]float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var lat, lon float64
		if _, err := fmt.Sscanf(sc.Text(), "%f,%f", &lat, &lon); err == nil {
			points = append(points, [2]float64{lat, lon})
		}
	}
	return points, sc.Err()
}

// Clear truncates the track file.
func (t *TrackLog) Clear() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.file.Truncate(0)
}

// Close closes the track file.
func (t *TrackLog) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.file.Close()
}
