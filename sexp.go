package gore

import (
	"encoding/binary"
	"errors"
	"math"
)

func parseSEXP(buf []byte, offset int) (interface{}, error) {

	//hasAttribute := buf[offset]&128 != 0
	//isLong := buf[offset]&64 != 0
	xt := expression(buf[offset] & 63)
	length := len(buf)
	offset = offset + 4

	if xt == XT_ARRAY_DOUBLE {
		noDoubles := length / 8
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
	return nil, errors.New("Unsupported expression type")
}
