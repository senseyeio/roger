package gore

type authType int

const (
	AT_plain authType = 1
	AT_crypt authType = 2
)

type command int

const (
	CMD_eval = 3
)

type typ int

const (
	DT_STRING typ = 4
	DT_SEXP   typ = 10
)
