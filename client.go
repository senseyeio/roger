package roger

import (
	"io"
	"net"
	"strconv"
)

// RClient is the main Roger interface allowing interaction with R.
type RClient interface {

	// Eval evaluates an R command synchronously returning the resulting object and any possible error. Creates a new session per command.
	Eval(command string) (interface{}, error)

	// Evaluate evaluates an R command asynchronously. The returned channel will resolve to a Packet once the command has completed. Creates a new session per command.
	Evaluate(command string) <-chan Packet

	// EvaluateSync evaluates an R command synchronously, resulting in a Packet. Creates a new session per command.
	EvaluateSync(command string) Packet

	// GetSession gets a session object which can be used to perform multiple commands in the same Rserve session.
	GetSession() (Session, error)
}

type roger struct {
	address  *net.TCPAddr
	user     string
	password string
}

// NewRClient creates a RClient which will run commands on the RServe server located at the provided host and port
func NewRClient(host string, port int64) (RClient, error) {
	return NewRClientWithAuth(host, port, "", "")
}

// NewRClientWithAuth creates a RClient with the specified credentials and RServe server details
func NewRClientWithAuth(host string, port int64, user, password string) (RClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}

	rClient := &roger{
		address:  addr,
		user:     user,
		password: password,
	}

	if _, err = rClient.Eval("'Test session connection'"); err != nil {
		return nil, err
	}
	return rClient, nil
}

func (r *roger) GetSession() (Session, error) {
	rwc, err := r.getReadWriteCloser()
	if err != nil {
		return nil, err
	}
	return newSession(rwc, r.user, r.password)
}

func (r *roger) getReadWriteCloser() (io.ReadWriteCloser, error) {
	connection, err := net.DialTCP("tcp", nil, r.address)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (r *roger) EvaluateSync(command string) Packet {
	sess, err := r.GetSession()
	if err != nil {
		return newErrorPacket(err)
	}
	defer sess.Close()
	packet := sess.SendCommand(command + "\n")
	return packet
}

func (r *roger) Evaluate(command string) <-chan Packet {
	out := make(chan Packet)
	go func() {
		out <- r.EvaluateSync(command)
		close(out)
	}()
	return out
}

func (r *roger) Eval(command string) (interface{}, error) {
	return r.EvaluateSync(command).GetResultObject()
}
