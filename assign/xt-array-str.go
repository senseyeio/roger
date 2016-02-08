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
	if rl > 0xfffff0 {
		rl += 4
	}
	rl += 4

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
		setHdrOffset(34, rl-8, rq, off)
		off += 8
	} else {
		off = sl + 8
		setHdrOffset(34, rl-4, rq, off)
		off += 4
	}

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
