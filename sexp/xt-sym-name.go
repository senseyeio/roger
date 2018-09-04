package sexp

func parseSymName(buf []byte, offset, end int) (interface{}, int, error) {
	endOfString := offset
	for buf[endOfString] != 0 && endOfString < end {
		endOfString = endOfString + 1
	}
	return string(buf[offset:endOfString]), end, nil
}
