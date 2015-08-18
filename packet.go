package roger

import (
	"errors"
	"strconv"

	"github.com/dareid/gore/sexp"
)

type Packet struct {
	cmd     int
	content []byte
	err     error
}

func newPacket(cmd int, content []byte) *Packet {
	return &Packet{
		cmd:     cmd,
		content: content,
	}
}

func newErrorPacket(err error) *Packet {
	return &Packet{
		err: err,
	}
}

func (p *Packet) IsError() bool {
	return p.err != nil || p.cmd&15 == 2
}

func (p *Packet) getStatusCode() int {
	return p.cmd >> 24 & 127
}

func (p *Packet) GetError() error {
	if p.err != nil {
		return p.err
	}
	return errors.New("Command error with status: " + strconv.Itoa(p.getStatusCode()))
}

func (p *Packet) GetResultObject() (interface{}, error) {
	if p.IsError() {
		return nil, p.GetError()
	}
	isSexp := p.content[0] == byte(dtSexp)
	if !isSexp {
		return nil, errors.New("Expected SEXP response")
	}
	return sexp.Parse(p.content[4:len(p.content)], 0)
}
