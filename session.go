package roger

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/senseyeio/roger/assign"
	"github.com/senseyeio/roger/constants"
)

// Session is an interface to send commands to an RServe session. Sessions must be closed after use.
type Session interface {

	// Sends a command to RServe which is evaluated synchronously resulting in a Packet.
	SendCommand(command string) Packet

	// Eval evaluates an R command synchronously returning the resulting object and any possible error. Unlike client.Eval, this does not start a new session.
	Eval(cmd string) (interface{}, error)

	// Assign value to a variable within the R session.
	Assign(symbol string, cont interface{}) error

	// Close closes a RServe session. Sessions must be closed after use.
	Close()
}

type authType int

const (
	atPlain authType = 1
	atCrypt authType = 2
)

type session struct {
	readWriteClose     io.ReadWriteCloser
	readWriter         *bufio.ReadWriter
	authReq            bool
	authType           authType
	key                string
	connected          bool
	rServeIDSig        string
	rServeProtocol     string
	rServeCommProtocol string

	user     string
	password string
}

func newSession(readWriteCloser io.ReadWriteCloser, user, password string) (*session, error) {
	buffRead := bufio.NewReader(readWriteCloser)
	buffWrite := bufio.NewWriter(readWriteCloser)
	sess := &session{
		readWriteClose: readWriteCloser,
		readWriter:     bufio.NewReadWriter(buffRead, buffWrite),
		user:           user,
		password:       password,
	}
	err := sess.handshake()
	return sess, err
}

func (s *session) Close() {
	s.connected = false
	s.readWriter = nil
	if s.readWriteClose != nil {
		s.readWriteClose.Close()
	}
	s.readWriteClose = nil
}

func (s *session) readNBytes(bytes int) []byte {
	ret := make([]byte, bytes)
	for v := 0; v < bytes; v++ {
		ret[v], _ = s.readWriter.ReadByte()
	}
	return ret
}

func (s *session) toCharset(str string) []byte {
	return []byte(str)
}

func (s *session) parseInitialMessage() error {
	isByteArrayJustZeros := func(bArr []byte) bool {
		return len(bytes.Replace(bArr, []byte{0}, []byte{}, -1)) == 0
	}
	rserveIDSigBytes := s.readNBytes(4)
	rServeProtocolBytes := s.readNBytes(4)
	rServeCommProtocolBytes := s.readNBytes(4)
	if isByteArrayJustZeros(rserveIDSigBytes) ||
		isByteArrayJustZeros(rServeProtocolBytes) ||
		isByteArrayJustZeros(rServeCommProtocolBytes) {
		return errors.New("Handshake failed - please check the connection details")
	}

	s.rServeIDSig = string(rserveIDSigBytes)
	s.rServeProtocol = string(rServeProtocolBytes)
	s.rServeCommProtocol = string(rServeCommProtocolBytes)
	for i := 12; i < 32; i += 4 {
		attr := s.readNBytes(4)
		attrString := string(attr)
		if attrString == "ARpt" && s.authReq == false {
			s.authReq = true
			s.authType = atPlain
		}
		if attrString == "ARuc" {
			s.authReq = true
			s.authType = atCrypt
		}
		if attrString[0] == 'K' {
			s.key = attrString[1:3]
		}
	}
	return nil
}

func (s *session) login() error {
	if s.authReq == false {
		return nil
	}
	if s.authReq == true && (s.user == "" || s.password == "") {
		return errors.New("Authentication is required and no credentials have been specified")
	}
	if s.key == "" {
		s.key = "rs"
	}
	cmd := s.user + "\n" + s.password
	if s.authType == atCrypt {
		cmd = s.user + "\n" + crypt(s.password, s.key)
	}

	packet := s.sendCommand(constants.CmdLogin, cmd)
	if packet.IsError() {
		_, err := packet.GetResultObject()
		return errors.New("Authentication failed: " + err.Error())
	}
	return nil
}

func (s *session) handshake() error {
	err := s.parseInitialMessage()
	if err != nil {
		return err
	}

	err = s.login()
	if err != nil {
		return err
	}

	s.connected = true

	if s.rServeCommProtocol != "QAP1" ||
		s.rServeIDSig != "Rsrv" ||
		s.rServeProtocol != "0103" {
		log.Println("The version of RServe installed is not officially supported. Please consider upgrading to the latest version of RServe.")
	}
	return nil
}

func (s *session) prepareStringCommand(cmd string) []byte {
	cmd = strings.Replace(cmd, "\r", "\n", -1) //avoid potential issue when loading external r script block
	rawCmdBytes := s.toCharset(cmd)
	requiredLength := len(rawCmdBytes) + 1
	//make sure length is divisible by 4
	if requiredLength&3 > 0 {
		requiredLength = (requiredLength & 0xfffffc) + 4
	}
	hdrLength := assign.GetHeaderLength(constants.DtString, requiredLength)
	cmdBytes := make([]byte, requiredLength+1+hdrLength)
	for i := 0; i < len(rawCmdBytes); i++ {
		cmdBytes[hdrLength+i] = rawCmdBytes[i]
	}
	assign.SetHdr(constants.DtString, requiredLength, cmdBytes)
	return cmdBytes
}

func (s *session) exeCommand(cmdType constants.Command, cmd string) {
	cmdBytes := s.prepareStringCommand(cmd)
	buf := new(bytes.Buffer)
	//command
	binary.Write(buf, binary.LittleEndian, int32(cmdType))
	//length of message (bits 0-31)
	binary.Write(buf, binary.LittleEndian, int32(len(cmdBytes)))
	//offset of message part
	binary.Write(buf, binary.LittleEndian, int32(0))
	// length of message (bits 32-63)
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, cmdBytes)

	s.readWriter.Write(buf.Bytes())
	s.readWriter.Flush()
}

func (s *session) readResponse() Packet {
	rep := binary.LittleEndian.Uint32(s.readNBytes(4))
	r1 := binary.LittleEndian.Uint32(s.readNBytes(4))
	s.readNBytes(8)

	if r1 <= 0 {
		return newPacket(int(rep), nil)
	}

	results := s.readNBytes(int(r1))
	return newPacket(int(rep), results)
}

func (s *session) sendCommand(cmdType constants.Command, cmd string) Packet {
	if s.connected == false && cmdType != constants.CmdLogin {
		return newErrorPacket(errors.New("Session was previously closed"))
	}
	s.exeCommand(cmdType, cmd)
	return s.readResponse()
}

func (s *session) SendCommand(cmd string) Packet {
	return s.sendCommand(constants.CmdEval, cmd)
}

func (s *session) Eval(cmd string) (interface{}, error) {
	return s.sendCommand(constants.CmdEval, cmd).GetResultObject()
}

func (s *session) request(cmdType constants.Command, cont []byte, offset int, length int) Packet {
	if cont != nil {
		if offset >= len(cont) {
			cont = nil
			length = 0
		} else if length > (len(cont) - offset) {
			length = len(cont) - offset
		}
	}
	if offset < 0 {
		offset = 0
	}
	if length < 0 {
		length = 0
	}

	contlen := 0
	if cont != nil {
		contlen = length
	}

	hdr := make([]byte, 16)
	assign.SetInt(int(cmdType), hdr, 0)
	assign.SetInt(contlen, hdr, 4)
	for i := 8; i < 16; i++ {
		hdr[i] = 0
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, hdr)
	binary.Write(buf, binary.LittleEndian, cont)

	s.readWriter.Write(buf.Bytes())
	s.readWriter.Flush()

	return s.readResponse()
}

func (s *session) Assign(symbol string, value interface{}) error {
	assignCommand, err := assign.Assign(symbol, value)
	if err != nil {
		return err
	}
	rp := s.request(constants.CmdSetSexp, assignCommand, 0, len(assignCommand))
	if rp != nil && rp.IsOk() {
		return nil
	}
	return errors.New("Assign failed")
}
