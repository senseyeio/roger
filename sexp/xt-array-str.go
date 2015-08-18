package sexp

func parseArrayString(buf []byte, offset int) (interface{}, error) {
	length := len(buf)
	noStrings := 0
	for ct := offset; ct < length; ct++ {
		if buf[ct] == 0 {
			noStrings = noStrings + 1
		}
	}
	stringArr := make([]string, noStrings, noStrings)

	startOfString := offset
	idx := 0
	for offset < length && idx < noStrings {
		if buf[offset] == 0 {
			endOfString := startOfString + (offset - startOfString)
			if buf[startOfString] == 0xff {
				if buf[startOfString+1] == 0 {
					stringArr[idx] = ""
				} else {
					stringArr[idx] = string(buf[startOfString+1 : endOfString-1])
				}
			} else {
				stringArr[idx] = string(buf[startOfString:endOfString])
			}

			idx = idx + 1
			startOfString = offset + 1
		}
		offset = offset + 1
	}
	return stringArr, nil
}
