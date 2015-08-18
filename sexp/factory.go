package sexp

import (
	"errors"
	"strconv"
)

type expression int

const (
	XT_ARRAY_DOUBLE expression = 33
	XT_ARRAY_STR    expression = 34
)

func Parse(buf []byte, offset int) (interface{}, error) {
	xt := expression(buf[offset] & 63)
	offset = offset + 4

	if xt == XT_ARRAY_DOUBLE {
		return parseArrayDouble(buf, offset)
	}
	if xt == XT_ARRAY_STR {
		return parseArrayString(buf, offset)
	}
	return nil, errors.New("Unsupported expression type: " + strconv.Itoa(int(xt)))

}
