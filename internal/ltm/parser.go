package ltm

import "errors"

var (
	ErrChecksum    = errors.New("ltm: checksum mismatch")
	ErrPayloadSize = errors.New("ltm: unexpected payload size")
	ErrUnknownFunc = errors.New("ltm: unknown function")
)

// Parser state constants.
const (
	stateIdle    = iota
	stateHeader  // got '$', expecting 'T'
	stateFunc    // got '$T', expecting function char
	statePayload // accumulating payload bytes
	stateCheck   // expecting checksum byte
)

// Parser is a push-model LTM frame parser implementing io.Writer.
type Parser struct {
	state   int
	fn      byte // current frame function
	size    int  // expected payload size
	buf     []byte
	Handler func(RawFrame)
	OnError func(error)
}

// NewParser creates a Parser that calls handler for each valid frame.
func NewParser(handler func(RawFrame), onError func(error)) *Parser {
	return &Parser{
		state:   stateIdle,
		buf:     make([]byte, 0, 16),
		Handler: handler,
		OnError: onError,
	}
}

// Write implements io.Writer.
func (p *Parser) Write(data []byte) (int, error) {
	for _, b := range data {
		p.feed(b)
	}
	return len(data), nil
}

func (p *Parser) feed(b byte) {
	switch p.state {
	case stateIdle:
		if b == Header1 {
			p.state = stateHeader
		}
	case stateHeader:
		if b == Header2 {
			p.state = stateFunc
		} else {
			p.state = stateIdle
			if b == Header1 {
				p.state = stateHeader
			}
		}
	case stateFunc:
		size, ok := PayloadSize[b]
		if !ok {
			if p.OnError != nil {
				p.OnError(ErrUnknownFunc)
			}
			p.state = stateIdle
			if b == Header1 {
				p.state = stateHeader
			}
			return
		}
		p.fn = b
		p.size = size
		p.buf = p.buf[:0]
		p.state = statePayload
	case statePayload:
		p.buf = append(p.buf, b)
		if len(p.buf) >= p.size {
			p.state = stateCheck
		}
	case stateCheck:
		expected := xorChecksum(p.buf)
		if b != expected {
			if p.OnError != nil {
				p.OnError(ErrChecksum)
			}
			// Resync: check if this byte is '$'
			p.state = stateIdle
			if b == Header1 {
				p.state = stateHeader
			}
			return
		}

		// Valid frame — copy payload so caller can hold reference
		payload := make([]byte, len(p.buf))
		copy(payload, p.buf)

		frame := RawFrame{
			Function: p.fn,
			Payload:  payload,
		}

		p.state = stateIdle
		if p.Handler != nil {
			p.Handler(frame)
		}
	}
}

// xorChecksum computes XOR of all payload bytes.
// INAV's LTM checksum does NOT include the function byte.
func xorChecksum(payload []byte) byte {
	crc := byte(0)
	for _, b := range payload {
		crc ^= b
	}
	return crc
}
