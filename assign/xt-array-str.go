package assign

import "github.com/senseyeio/roger/constants"

func assignStrArray(symbol string, value []string) ([]byte, error) {
	rl := 0
	i := 0
	for i < len(value) {
		b := []byte(value[i])
		if len(b) > 0 {
			rl += len(b)
		}
		rl++
		i++
	}
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
	sextHeader := GetHeaderLength(constants.DataType(constants.XtStringArray), rl)
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
	setHdrOffset(constants.DataType(constants.XtStringArray), rl, rq, off)
	off += sextHeader

	i = 0
	io := off
	for i < len(value) {
		b := []byte(value[i])
		if len(b) > 0 {
			copy(rq[io:io+len(b)], b[:])
			io += len(b)
		}
		rq[io] = 0
		io++
		i++
	}
	i = io - off
	for (i & 3) != 0 {
		rq[io] = 1
		io++
		i++
	}

	return rq, nil

}
