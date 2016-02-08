package assign

import "github.com/senseyeio/roger/constants"

func assignStr(symbol string, value string) ([]byte, error) {
	symn := []byte(symbol)
	ctn := []byte(value)

	sl := len(symn) + 1
	cl := len(ctn) + 1

	if (sl & 3) > 0 {
		sl = (sl & 0xfffffc) + 4
	}
	if (cl & 3) > 0 {
		cl = (cl & 0xfffffc) + 4
	}

	rq := make([]byte, sl+4+cl+4)

	ic := 0
	for ; ic < len(symn); ic++ {
		rq[ic+4] = symn[ic]
	}
	for ic < sl {
		rq[ic+4] = 0
		ic++
	}
	for ic = 0; ic < len(ctn); ic++ {
		rq[ic+sl+8] = ctn[ic]
	}
	for ic < cl {
		rq[ic+sl+8] = 0
		ic++
	}

	setHdrOffset(constants.DtString, sl, rq, 0)
	setHdrOffset(constants.DtString, cl, rq, sl+4)

	return rq, nil
}
