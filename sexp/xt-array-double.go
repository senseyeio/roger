package sexp

import (
	"encoding/binary"
	"math"
)

func parseArrayDouble(buf []byte, offset int) (interface{}, error) {
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
	return doubleArr, nil
}
