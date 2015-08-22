package sexp

import "encoding/binary"

func parseInt(buf []byte, offset int, end int) (interface{}, int, error) {
	bits := binary.LittleEndian.Uint32(buf[offset : offset+4])
	return int32(bits), offset + 4, nil
}
