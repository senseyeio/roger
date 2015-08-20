package roger

import (
	"errors"
	"strconv"

	"github.com/senseyeio/roger/sexp"
)

// Packet is the interface satisfied by objects returned from a R command.
// It contains either the resulting object or an error.
type Packet interface {
	GetResultObject() (interface{}, error)
	IsError() bool
	GetError() error
}

type packet struct {
	cmd     int
	content []byte
	err     error
}

func newPacket(cmd int, content []byte) Packet {
	return &packet{
		cmd:     cmd,
		content: content,
	}
}

func newErrorPacket(err error) Packet {
	return &packet{
		err: err,
	}
}

// IsError returns a boolean defining whether the Packet contains an error.
func (p *packet) IsError() bool {
	return p.err != nil || p.cmd&15 == 2
}

func (p *packet) getStatusCode() int {
	return p.cmd >> 24 & 127
}

// GetError returns an error if the Packet contains an error. If not it returns nil.
func (p *packet) GetError() error {
	if p.IsError() == false {
		return nil
	}
	if p.err != nil {
		return p.err
	}
	return errors.New("Command error with status: " + strconv.Itoa(p.getStatusCode()))
}

// GetResultObject will parse the packet's contents, returning a go interface{}.
// If the Packet contains an error this will be returned.
// If conversion fails, an error will be returned.
func (p *packet) GetResultObject() (interface{}, error) {
	if p.IsError() {
		return nil, p.GetError()
	}
	if len(p.content) == 0 {
		return nil, errors.New("Command failed for an unknown reason")
	}
	isSexp := p.content[0] == byte(dtSexp)
	if !isSexp {
		return nil, errors.New("Expected SEXP response")
	}
	return sexp.Parse(p.content[4:len(p.content)], 0)
}
