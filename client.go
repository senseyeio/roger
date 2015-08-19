package roger

import (
	"io"
	"net"
	"strconv"
)

// RClient is the main Roger interface allowing interaction with R.
type RClient interface {
	Evaluate(command string) <-chan *Packet
	EvaluateSync(command string) *Packet
	getReadWriteCloser() (io.ReadWriteCloser, error)
}

type roger struct {
	address  *net.TCPAddr
	user     string
	password string
}

// NewRClient creates a RClient which will run commands on the RServe server located at the provided host and port
func NewRClient(host string, port int64) (RClient, error) {
	return newRClientWithAuth(host, port, "", "")
}

func newRClientWithAuth(host string, port int64, user, password string) (RClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}

	return &roger{
		address:  addr,
		user:     user,
		password: password,
	}, nil
}

// EvaluateSync evaluates a R command synchronously.
func (r *roger) EvaluateSync(command string) *Packet {
	sess, err := newSession(r)
	if err != nil {
		return newErrorPacket(err)
	}
	packet := sess.sendCommand(command + "\n")
	sess.close()
	return packet
}

// Evaluate evaluates a R command asynchronously. The returned channel will resolve to a Packet once the command has completed.
func (r *roger) Evaluate(command string) <-chan *Packet {
	out := make(chan *Packet)
	go func() {
		out <- r.EvaluateSync(command)
		close(out)
	}()
	return out
}

func (r *roger) getReadWriteCloser() (io.ReadWriteCloser, error) {
	connection, err := net.DialTCP("tcp", nil, r.address)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
