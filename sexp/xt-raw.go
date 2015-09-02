package sexp

import "encoding/binary"

func parseRaw(buf []byte, offset int, end int) (interface{}, int, error) {
	as := binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset = offset + 4
	rawBuffer := buf[offset : offset+int(as)]
	return rawBuffer, end, nil
}
