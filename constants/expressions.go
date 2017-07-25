package constants

type ExpressionType int

const (
	XtNull         ExpressionType = 0  /* P data: [0] */
	XtInt          ExpressionType = 1  /* P data: [4]int */
	XtDouble       ExpressionType = 2  /* -  data: [8]double */
	XtString       ExpressionType = 3  /* P data: [n]char null-term. strg. */
	XtLang         ExpressionType = 4  /* - */
	XtS4           ExpressionType = 7  /* P data: [0] */
	XtVector       ExpressionType = 16 /* P data: [?]REXP,REXP,.. */
	XtClos         ExpressionType = 18 /* P X formals, X body  (closure; since 0.1-5) */
	XtSymName      ExpressionType = 19 /* s same as XtStr (since 0.5) */
	XtListNoTag    ExpressionType = 20 /* s same as XtVector (since 0.5) */
	XtListTag      ExpressionType = 21 /* P X tag, X val, Y tag, Y val, ... (since 0.5) */
	XtLangNoTag    ExpressionType = 22 /* s same as XtListNoTag (since 0.5) */
	XtLangTag      ExpressionType = 23 /* s same as XtListTag (since 0.5) */
	XtExpVector    ExpressionType = 26 /* s same as XtVector (since 0.5) */
	XtIntArray     ExpressionType = 32 /* P data: [n*4]int,int,.. */
	XtDoubleArray  ExpressionType = 33 /* P data: [n*8]double,double,.. */
	XtStringArray  ExpressionType = 34 /* P data: string,string,.. (string=byte,byte,...,0) padded with '\01' */
	XtBoolArray    ExpressionType = 36 /* P data: int(n),byte,byte,... */
	XtRaw          ExpressionType = 37 /* P data: int(n),byte,byte,... */
	XtComplexArray ExpressionType = 38 /* P data: [n*16]double,double,... (Re,Im,Re,Im,...) */
	XtUnknown      ExpressionType = 48 /* P data: [4]int - SEXP type (as from TYPEOF(x)) */
)
