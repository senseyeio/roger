package roger

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"errors"
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
}

func newSession(client RClient) (*session, error) {
	readWriteCloser, err := client.getReadWriteCloser()
	if err != nil {
		return nil, err
	}
	buffRead := bufio.NewReader(readWriteCloser)
	buffWrite := bufio.NewWriter(readWriteCloser)
	sess := &session{
		readWriteClose: readWriteCloser,
		readWriter:     bufio.NewReadWriter(buffRead, buffWrite),
	}
	err = sess.handshake()
	return sess, err
}

func (s *session) close() {
	s.connected = false
	s.readWriter = nil
	s.readWriteClose.Close()
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

func (s *session) handshake() error {
	s.rServeIDSig = string(s.readNBytes(4))
	s.rServeProtocol = string(s.readNBytes(4))
	s.rServeCommProtocol = string(s.readNBytes(4))
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
	s.connected = true
	if s.rServeCommProtocol == "" ||
		s.rServeIDSig == "" ||
		s.rServeProtocol == "" {
		return errors.New("Handshake failed");
	}
	if s.rServeCommProtocol != "QAP1" ||
		s.rServeIDSig != "Rsrv" ||
		s.rServeProtocol != "0103" {
		log.Println("The version of RServe installed is not officially supported. Please consider upgrading to the latest version of RServe.")
	}
	return nil
}

func (s *session) setHdr(valueType dataType, valueLength int, buf []byte) {
	buf[0] = byte(valueType)
	buf[1] = byte(valueLength & 255)
	buf[2] = byte((valueLength & 0xff00) >> 8)
	buf[3] = byte((valueLength & 0xff0000) >> 16)
}

func (s *session) prepareStringCommand(cmd string) []byte {
	rawCmdBytes := s.toCharset(cmd)
	requiredLength := len(rawCmdBytes) + 1
	//make sure length is divisible by 4
	if requiredLength&3 > 0 {
		requiredLength = (requiredLength & 0xfffffc) + 4
	}
	cmdBytes := make([]byte, requiredLength+5)
	for i := 0; i < len(rawCmdBytes); i++ {
		cmdBytes[4+i] = rawCmdBytes[i]
	}
	s.setHdr(dtString, requiredLength, cmdBytes)
	return cmdBytes
}

func (s *session) sendCommand(cmd string) *Packet {
	cmdBytes := s.prepareStringCommand(cmd)
	buf := new(bytes.Buffer)
	//command
	binary.Write(buf, binary.LittleEndian, int32(cmdEval))
	//length of message (bits 0-31)
	binary.Write(buf, binary.LittleEndian, int32(len(cmdBytes)))
	//offset of message part
	binary.Write(buf, binary.LittleEndian, int32(0))
	// length of message (bits 32-63)
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, cmdBytes)

	s.readWriter.Write(buf.Bytes())
	s.readWriter.Flush()

	rep := binary.LittleEndian.Uint32(s.readNBytes(4))
	r1 := binary.LittleEndian.Uint32(s.readNBytes(4))
	s.readNBytes(8)

	if r1 <= 0 {
		return newPacket(int(rep), nil)
	}

	results := s.readNBytes(int(r1))
	return newPacket(int(rep), results)
}
