package sexp

import (
	"encoding/binary"
	"math"
)

func parseDoubleArray(buf []byte, offset, end int) (interface{}, int, error) {
	length := end - offset
	doubleArr := make([]float64, 0, length/8)
	for ; offset < end; offset += 8 {
		bits := binary.LittleEndian.Uint64(buf[offset : offset+8])
		doubleArr = append(doubleArr, math.Float64frombits(bits))
	}
	if len(doubleArr) == 1 {
		return doubleArr[0], offset, nil
	}
	return doubleArr, offset, nil
}
