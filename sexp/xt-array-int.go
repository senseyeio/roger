package sexp

import "encoding/binary"

func parseIntArray(buf []byte, offset, end int) (interface{}, int, error) {
	length := end - offset
	noInts := length / 4
	intArr := make([]int32, noInts, noInts)
	for ct := 0; ct < noInts; ct++ {
		start := offset
		end := start + 4
		bits := binary.LittleEndian.Uint32(buf[start:end])
		intArr[ct] = int32(bits)
		offset += 4
	}
	if len(intArr) == 1 {
		return intArr[0], offset, nil
	}
	return intArr, offset, nil
}
