package sexp

import (
	"encoding/binary"
	"math"
)

func parseComplexArray(buf []byte, offset, end int) (interface{}, int, error) {
	length := end - offset
	cArr := make([]complex128, 0, length/16)
	for ; offset < end; offset += 16 {
		bitsReal := binary.LittleEndian.Uint64(buf[offset : offset+8])
		bitsImag := binary.LittleEndian.Uint64(buf[offset+8 : offset+16])
		cArr = append(cArr, complex(math.Float64frombits(bitsReal), math.Float64frombits(bitsImag)))
	}
	if len(cArr) == 1 {
		return cArr[0], offset, nil
	}
	return cArr, offset, nil
}
