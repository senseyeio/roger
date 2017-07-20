package assign

import (
	"math"

	"github.com/senseyeio/roger/constants"
)

func assignDouble(symbol string, value float64) ([]byte, error) {
	rl := 8

	symn := []byte(symbol)
	sl := len(symn) + 1
	if (sl & 3) > 0 {
		sl = (sl & 0xfffffc) + 4
	}

	var rq []byte

	shl := GetHeaderLength(constants.DtString, sl)
	sextHeader := GetHeaderLength(constants.DataType(constants.XtDouble), rl)
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
	setHdrOffset(constants.DataType(constants.XtDouble), rl, rq, off)
	off += sextHeader

	setLong(int64(math.Float64bits(value)), rq, off)

	return rq, nil
}
