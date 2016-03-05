package sexp

import (
	"encoding/binary"

	"github.com/senseyeio/roger/types"
)

func convertRBoolArray(rbools []types.RBool) []bool {
	converted := make([]bool, len(rbools))
	for i, rb := range rbools {
		nativeBool, _ := rb.ToBoolean()
		converted[i] = nativeBool
	}
	return converted
}

func parseBoolArray(buf []byte, offset int, end int) (interface{}, int, error) {
	boolArrayLen := binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset = offset + 4

	rBoolArr := make([]types.RBool, 0, boolArrayLen)
	for ct := uint32(0); ct < boolArrayLen; ct++ {
		b := buf[offset]
		if b == 1 {
			rBoolArr = append(rBoolArr, types.TRUE)
		} else if b == 0 {
			rBoolArr = append(rBoolArr, types.FALSE)
		} else {
			rBoolArr = append(rBoolArr, types.NA)
		}
		offset = offset + 1
	}

	//boolArr := convertRBoolArray(rBoolArr)

	if len(rBoolArr) == 1 {
		return rBoolArr[0], end, nil
	}

	return rBoolArr, end, nil
}
