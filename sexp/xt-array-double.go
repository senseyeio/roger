package sexp

import (
	"encoding/binary"
	"math"
)

func parseDoubleArray(buf []byte, offset int) (interface{}, int, error) {
	length := len(buf)
	noDoubles := (length - offset) / 8
	doubleArr := make([]float64, noDoubles, noDoubles)
	for ct := 0; ct < noDoubles; ct++ {
		start := offset
		end := start + 8
		bits := binary.LittleEndian.Uint64(buf[start:end])
		doubleArr[ct] = math.Float64frombits(bits)
		offset += 8
	}
	if len(doubleArr) == 1 {
		return doubleArr[0], offset, nil
	}
	return doubleArr, offset, nil
}
