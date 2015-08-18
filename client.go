package gore

import (
	"io"
	"net"
	"strconv"
)

type RClient interface {
	Evaluate(command string) <-chan *Packet
	EvaluateSync(command string) *Packet
	getReadWriteCloser() (io.ReadWriteCloser, error)
}

type gore struct {
	address  *net.TCPAddr
	user     string
	password string
}

func NewRClient(host string, port int64) (RClient, error) {
	return NewRClientWithAuth(host, port, "", "")
}

func NewRClientWithAuth(host string, port int64, user, password string) (RClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}

	return &gore{
		address:  addr,
		user:     user,
		password: password,
	}, nil
}

func (gore *gore) EvaluateSync(command string) *Packet {
	sess, err := newSession(gore)
	if err != nil {
		return newErrorPacket(err)
	}
	packet := sess.sendCommand(command + "\n")
	sess.close()
	return packet
}

func (gore *gore) Evaluate(command string) <-chan *Packet {
	out := make(chan *Packet)
	go func() {
		out <- gore.EvaluateSync(command)
		close(out)
	}()
	return out
}

func (gore *gore) getReadWriteCloser() (io.ReadWriteCloser, error) {
	connection, err := net.DialTCP("tcp", nil, gore.address)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
