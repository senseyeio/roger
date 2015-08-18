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
	strIdx := 0
	for offset < length {
		if buf[offset] == 0 {
			stringLength := offset - startOfString
			if buf[startOfString] == 0xff {
				if buf[startOfString+1] == 0 {
					stringArr[strIdx] = ""
				} else {
					stringArr[strIdx] = string(buf[startOfString+1 : startOfString+(stringLength-1)])
				}
			} else {
				stringArr[strIdx] = string(buf[startOfString : startOfString+(stringLength)])
			}

			strIdx = strIdx + 1
			startOfString = offset + 1
		}
		offset = offset + 1
	}
	return stringArr, nil
}
