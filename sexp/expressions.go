package sexp

type expressionType int

const (
	xtNull         expressionType = 0
	xtInt          expressionType = 1
	xtVector       expressionType = 16
	xtSymName      expressionType = 19
	xtListTag      expressionType = 21
	xtIntArray     expressionType = 32
	xtDoubleArray  expressionType = 33
	xtStringArray  expressionType = 34
	xtBoolArray    expressionType = 36
	xtRaw          expressionType = 37
	xtComplexArray expressionType = 38
)
