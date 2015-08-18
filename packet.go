package gore

import (
	"errors"

	"github.com/dareid/gore/sexp"
)

type Packet struct {
	cmd     int
	content []byte
}

func NewPacket(cmd int, content []byte) *Packet {
	return &Packet{
		cmd:     cmd,
		content: content,
	}
}

func (p *Packet) IsOK() bool {
	return p.cmd&15 == 1
}

func (p *Packet) IsError() bool {
	return p.cmd&15 == 2
}

func (p *Packet) GetStatusCode() int {
	return p.cmd >> 24 & 127
}

func (p *Packet) GetResultObject() (interface{}, error) {
	isSexp := p.content[0] == byte(DT_SEXP)
	if !isSexp {
		return nil, errors.New("Expected SEXP response")
	}
	return sexp.Parse(p.content[4:len(p.content)], 0)
}
