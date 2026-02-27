package serial

import (
	"strings"
	"testing"
)

func TestOpen_InvalidPort(t *testing.T) {
	_, err := Open(Config{Name: "/dev/nonexistent-port-12345", Baud: 115200})
	if err == nil {
		t.Fatal("expected error for nonexistent port")
	}
	if !strings.Contains(err.Error(), "open serial") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{Name: "/dev/ttyUSB0", Baud: 19200}
	if cfg.Name != "/dev/ttyUSB0" {
		t.Errorf("Name=%s, want /dev/ttyUSB0", cfg.Name)
	}
	if cfg.Baud != 19200 {
		t.Errorf("Baud=%d, want 19200", cfg.Baud)
	}
}
