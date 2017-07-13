package assign

import (
	"github.com/senseyeio/roger/constants"
)

func assignStr(symbol string, value string) ([]byte, error) {
	symn := []byte(symbol)
	ctn := []byte(value)

	sl := len(symn) + 1
	cl := len(ctn) + 1

	shl := GetHeaderLength(constants.DtString, sl)
	chl := GetHeaderLength(constants.DtString, cl)

	rq := make([]byte, sl+shl+cl+chl)

	ic := 0
	for ; ic < len(symn); ic++ {
		rq[ic+shl] = symn[ic]
	}
	for ic < sl {
		rq[ic+shl] = 0
		ic++
	}
	for ic = 0; ic < len(ctn); ic++ {
		rq[ic+sl+shl+chl] = ctn[ic]
	}
	for ic < cl {
		rq[ic+sl+shl+chl] = 0
		ic++
	}

	setHdrOffset(constants.DtString, sl, rq, 0)
	setHdrOffset(constants.DtString, cl, rq, sl+shl)

	return rq, nil
}
