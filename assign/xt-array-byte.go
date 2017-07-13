package assign

import "github.com/senseyeio/roger/constants"

func assignByteArray(symbol string, value []byte) ([]byte, error) {
	rl := len(value) + 4
	if (rl & 3) > 0 {
		rl = rl - (rl & 3) + 4
	}

	symn := []byte(symbol)
	sl := len(symn) + 1
	if (sl & 3) > 0 {
		sl = (sl & 0xfffffc) + 4
	}

	var rq []byte

	shl := GetHeaderLength(constants.DtString, sl)
	sextHeader := GetHeaderLength(constants.DataType(constants.XtRaw), rl)
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
	setHdrOffset(constants.DtSexp, rl, rq, sl+shl)

	off := sl + shl + rhl
	setHdrOffset(constants.DataType(constants.XtRaw), rl+sextHeader, rq, off)
	off += sextHeader

	SetInt(len(value), rq, off)
	off += 4
	copy(rq[off:off+len(value)], value)

	return rq, nil
}
