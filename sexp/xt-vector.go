package sexp

import (
	"errors"
	"log"
)

func parseVectorAttr(attr interface{}, vectorArr []interface{}, offset int) (interface{}, int, error) {
	attrMap, ok := attr.(map[string]interface{})
	if !ok {
		return vectorArr, offset, nil
	}
	names, ok := attrMap["n"].([]string)
	if !ok {
		return vectorArr, offset, nil
	}
	if len(names) != len(vectorArr) {
		return nil, offset, errors.New("Vector name and value quantity mismatch")
	}

	ret := map[string]interface{}{}
	for i := range names {
		ret[names[i]] = vectorArr[i]
	}

	return ret, offset, nil
}

func parseVector(attr interface{}, buf []byte, offset, end int) (interface{}, int, error) {
	var vectorArr []interface{}

	for offset < end {
		vectorEntry, newOffset, err := parseReturningOffset(buf, offset)
		offset = newOffset
		if err != nil {
			vectorArr = append(vectorArr, err)
		} else {
			vectorArr = append(vectorArr, vectorEntry)
		}
	}
	if offset != end {
		log.Println("Warning: vector size mismatch")
		offset = end
	}
	return parseVectorAttr(attr, vectorArr, offset)
}
