package sexp

import (
	"errors"
	"strconv"
)

// Parse converts a byte array containing R SEXP to a golang object.
// This can be converted to native golang types.
func Parse(buf []byte, offset int) (interface{}, error) {
	obj, _, err := parseReturningOffset(buf, offset)
	return obj, err
}

func parseReturningOffset(buf []byte, offset int) (interface{}, int, error) {
	xt := expressionType(buf[offset] & 63)
	offset = offset + 4

	if xt == xtDoubleArray {
		return parseDoubleArray(buf, offset)
	}
	if xt == xtStringArray {
		return parseStringArray(buf, offset)
	}
	if xt == xtIntArray {
		return parseIntArray(buf, offset)
	}
	if xt == xtBoolArray {
		return parseBoolArray(buf, offset)
	}
	return nil, offset, errors.New("Unsupported expression type: " + strconv.Itoa(int(xt)))
}
