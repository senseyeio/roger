// Package sexp parses R s expression trees into native go objects
package sexp

import (
	"encoding/binary"
	"errors"
	"strconv"
)

// Parse converts a byte array containing R SEXP to a golang object.
// This can be converted to native golang types.
func Parse(buf []byte, offset int) (interface{}, error) {
	obj, _, err := parseReturningOffset(buf, offset)
	return obj, err
}

func getLength(buf []byte, offset int, isLong bool) int {
	if isLong {
		return int(binary.LittleEndian.Uint32(buf[offset+1 : offset+4]))
	}
	return int(uint32(buf[offset+1]) | (uint32(buf[offset+2]) << 8) | (uint32(buf[offset+3]) << 16))
}

func parseReturningOffset(buf []byte, offset int) (interface{}, int, error) {
	isLong := ((buf[offset] & 64) != 0)
	length := getLength(buf, offset, isLong)
	xt := expressionType(buf[offset] & 63)

	hasAtt := ((buf[offset] & 128) != 0)

	offset = offset + 4
	if isLong {
		offset = offset + 4
	}
	end := offset + length

	var attr interface{}
	if hasAtt {
		var err error
		attr, offset, err = parseReturningOffset(buf, offset)
		if err != nil {
			return nil, offset, err
		}
	}

	if xt == xtNull {
		return nil, offset, nil
	}
	if xt == xtInt {
		return parseInt(buf, offset, end)
	}
	if xt == xtSymName {
		return parseSymName(buf, offset, end)
	}
	if xt == xtDoubleArray {
		return parseDoubleArray(buf, offset, end)
	}
	if xt == xtStringArray {
		return parseStringArray(buf, offset, end)
	}
	if xt == xtIntArray {
		return parseIntArray(buf, offset, end)
	}
	if xt == xtBoolArray {
		return parseBoolArray(buf, offset, end)
	}
	if xt == xtVector {
		return parseVector(attr, buf, offset, end)
	}
	if xt == xtListTag {
		return parseListTag(buf, offset, end)
	}
	return nil, offset, errors.New("Unsupported expression type: " + strconv.Itoa(int(xt)))
}
