package sexp

import (
	"errors"
	"log"
)

func parseListTag(buf []byte, offset, end int) (interface{}, int, error) {
	list := map[string]interface{}{}

	for offset < end {
		var left, right interface{}
		var err error
		left, offset, err = parseReturningOffset(buf, offset)
		if err != nil {
			return nil, offset, err
		}
		right, offset, err = parseReturningOffset(buf, offset)
		if err != nil {
			return nil, offset, err
		}
		rightAsString, ok := right.(string)
		if !ok {
			return nil, offset, errors.New("Expecting xt-list-tag to have string tag")
		}
		list[rightAsString] = left
	}
	if offset != end {
		log.Println("Warning: List length mismatch")
	}
	return list, offset, nil
}
