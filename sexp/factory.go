package sexp

import (
	"errors"
	"strconv"
)

func Parse(buf []byte, offset int) (interface{}, error) {
	xt := expressionType(buf[offset] & 63)
	offset = offset + 4

	if xt == xtDoubleArray {
		return parseDoubleArray(buf, offset)
	}
	if xt == xtStringArray {
		return parseStringArray(buf, offset)
	}
	return nil, errors.New("Unsupported expression type: " + strconv.Itoa(int(xt)))

}
