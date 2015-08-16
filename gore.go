package gore

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type R interface {
	Evaluate(command string) *Packet
}

type gore struct {
	conn       *net.TCPConn
	readWriter *bufio.ReadWriter
	authReq    bool
	authType   authType
	key        string
	connected  bool
}

func NewRClient(host string, port int64, user, password string) (R, error) {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}
	connection, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	read := bufio.NewReader(connection)
	write := bufio.NewWriter(connection)

	g := &gore{
		conn:       connection,
		readWriter: bufio.NewReadWriter(read, write),
	}
	g.handshake()
	return g, nil
}

func (gore *gore) readNBytes(bytes int) []byte {
	ret := make([]byte, bytes)
	for v := 0; v < bytes; v++ {
		ret[v], _ = gore.readWriter.ReadByte()
	}
	return ret
}

func (gore *gore) skipNBytes(bytes int) {
	for v := 0; v < bytes; v++ {
		gore.readWriter.ReadByte()
	}
}

func (gore *gore) toCharset(str string) []byte {
	return []byte(str)
}

func (gore *gore) handshake() {
	rServeIDSig := gore.readNBytes(4)
	fmt.Println(string(rServeIDSig))
	rServeProtocol := gore.readNBytes(4)
	fmt.Println(string(rServeProtocol))
	rServeCommProtocol := gore.readNBytes(4)
	fmt.Println(string(rServeCommProtocol))
	for i := 12; i < 32; i += 4 {
		attr := gore.readNBytes(4)
		attrString := string(attr)
		fmt.Println(attrString)
		if attrString == "ARpt" && gore.authReq == false {
			gore.authReq = true
			gore.authType = AT_plain
		}
		if attrString == "ARuc" {
			gore.authReq = true
			gore.authType = AT_crypt
		}
		if attrString[0] == 'K' {
			gore.key = attrString[1:3]
		}
	}
	gore.connected = true
}

func (gore *gore) setHdr(valueType typ, valueLength int, buf []byte) {
	buf[0] = byte(valueType)
	buf[1] = byte(valueLength & 255)
	buf[2] = byte((valueLength & 0xff00) >> 8)
	buf[3] = byte((valueLength & 0xff0000) >> 16)
}

func (gore *gore) prepareStringCommand(cmd string) []byte {
	rawCmdBytes := gore.toCharset(cmd)
	requiredLength := len(rawCmdBytes) + 1
	//make sure length is divisible by 4
	if requiredLength&3 > 0 {
		requiredLength = (requiredLength & 0xfffffc) + 4
	}
	cmdBytes := make([]byte, requiredLength+5)
	for i := 0; i < len(rawCmdBytes); i++ {
		cmdBytes[4+i] = rawCmdBytes[i]
	}
	gore.setHdr(DT_STRING, requiredLength, cmdBytes)
	return cmdBytes
}

func (gore *gore) sendCommand(cmd string) *Packet {
	cmdBytes := gore.prepareStringCommand(cmd)
	buf := new(bytes.Buffer)
	//command
	binary.Write(buf, binary.LittleEndian, int32(CMD_eval))
	//length of message (bits 0-31)
	binary.Write(buf, binary.LittleEndian, int32(len(cmdBytes)))
	//offset of message part
	binary.Write(buf, binary.LittleEndian, int32(0))
	// length of message (bits 32-63)
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, cmdBytes)

	gore.readWriter.Write(buf.Bytes())
	gore.readWriter.Flush()

	rep := binary.LittleEndian.Uint32(gore.readNBytes(4))
	r1 := binary.LittleEndian.Uint32(gore.readNBytes(4))
	gore.skipNBytes(8)

	if r1 <= 0 {
		return nil
	}

	results := gore.readNBytes(int(r1))
	return NewPacket(int(rep), results)
}

func (gore *gore) Evaluate(command string) *Packet {
	return gore.sendCommand(command + "\n")
}
