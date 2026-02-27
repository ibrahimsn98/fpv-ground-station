package serial

import (
	"fmt"
	"time"

	goserial "go.bug.st/serial"
)

// Config holds serial port configuration.
type Config struct {
	Name string
	Baud int
}

// Port wraps a serial port as an io.ReadWriteCloser.
type Port struct {
	port goserial.Port
}

// Open opens a serial port with 8N1, no flow control, and a read timeout.
func Open(cfg Config) (*Port, error) {
	mode := &goserial.Mode{
		BaudRate: cfg.Baud,
		DataBits: 8,
		Parity:   goserial.NoParity,
		StopBits: goserial.OneStopBit,
	}

	p, err := goserial.Open(cfg.Name, mode)
	if err != nil {
		return nil, fmt.Errorf("open serial %s: %w", cfg.Name, err)
	}

	if err := p.SetReadTimeout(200 * time.Millisecond); err != nil {
		p.Close()
		return nil, fmt.Errorf("set read timeout: %w", err)
	}

	return &Port{port: p}, nil
}

// Read implements io.Reader.
func (p *Port) Read(buf []byte) (int, error) {
	return p.port.Read(buf)
}

// Write implements io.Writer.
func (p *Port) Write(buf []byte) (int, error) {
	return p.port.Write(buf)
}

// Close implements io.Closer.
func (p *Port) Close() error {
	return p.port.Close()
}

// Drain flushes the output buffer, ensuring all written bytes are transmitted.
func (p *Port) Drain() error {
	return p.port.Drain()
}

// SetReadTimeout adjusts the read timeout on the underlying serial port.
func (p *Port) SetReadTimeout(d time.Duration) error {
	return p.port.SetReadTimeout(d)
}

// ResetInputBuffer discards any data in the input buffer.
func (p *Port) ResetInputBuffer() error {
	return p.port.ResetInputBuffer()
}
