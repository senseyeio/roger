package assign

import "github.com/senseyeio/roger/constants"

func assignIntArray(symbol string, value []int32) ([]byte, error) {
	rl := len(value) * 4

	symn := []byte(symbol)
	sl := len(symn) + 1
	if (sl & 3) > 0 {
		sl = (sl & 0xfffffc) + 4
	}

	var rq []byte

	shl := GetHeaderLength(constants.DtString, sl)
	sextHeader := GetHeaderLength(constants.DataType(constants.XtIntArray), rl)
	rhl := GetHeaderLength(constants.DtSexp, rl+sextHeader)

	rq = make([]byte, sl+rl+shl+rhl+sextHeader)

	ic := 0
	for ; ic < len(symn); ic++ {
		rq[ic+shl] = symn[ic]
	}
	for ic < sl {
		rq[ic+shl] = 0
		ic++
	}

	setHdrOffset(constants.DtString, sl, rq, 0)
	setHdrOffset(constants.DtSexp, rl+sextHeader, rq, sl+shl)

	off := sl + shl + rhl
	setHdrOffset(constants.DataType(constants.XtIntArray), rl, rq, off)
	off += sextHeader

	i := 0
	io := off
	for i < len(value) {
		SetInt(int(value[i]), rq, io)
		i++
		io += 4
	}

	return rq, nil
}
