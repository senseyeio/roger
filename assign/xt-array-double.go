package assign

import (
	"math"

	"github.com/senseyeio/roger/constants"
)

func assignDoubleArray(symbol string, value []float64) ([]byte, error) {
	rl := len(value)*8 + 4
	if rl > 0xfffff0 {
		rl += 4
	}
	symn := []byte(symbol)
	sl := len(symn) + 1
	if (sl & 3) > 0 {
		sl = (sl & 0xfffffc) + 4
	}

	var rq []byte

	if rl > 0xfffff0 {
		rq = make([]byte, sl+rl+12)
	} else {
		rq = make([]byte, sl+rl+8)
	}

	ic := 0
	for ; ic < len(symn); ic++ {
		rq[ic+4] = symn[ic]
	}
	for ic < sl {
		rq[ic+4] = 0
		ic++
	}

	setHdrOffset(constants.DtString, sl, rq, 0)
	setHdrOffset(constants.DtSexp, rl, rq, sl+4)

	var off int
	if rl > 0xfffff0 {
		off = sl + 12
		setHdrOffset(33, rl-8, rq, off)
		off += 8
	} else {
		off = sl + 8
		setHdrOffset(33, rl-4, rq, off)
		off += 4
	}

	i := 0
	io := off
	for i < len(value) {
		setLong(int64(math.Float64bits(value[i])), rq, io)
		i++
		io += 8
	}

	return rq, nil
}
