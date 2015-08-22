package sexp

import (
	"errors"
	"log"
)

func parseListTag(buf []byte, offset, end int) (interface{}, int, error) {
	list := map[string]interface{}{}

	for offset < end {
		var left, right interface{}
		var err error
		left, offset, err = parseReturningOffset(buf, offset)
		if err != nil {
			return nil, offset, err
		}
		right, offset, err = parseReturningOffset(buf, offset)
		if err != nil {
			return nil, offset, err
		}
		rightAsString, ok := right.(string)
		if !ok {
			return nil, offset, errors.New("Expecting xt-list-tag to have string tag")
		}
		list[rightAsString] = left
	}
	if offset != end {
		log.Println("Warning: List length mismatch")
	}
	return list, offset, nil
}

/*
if (xt==XT_LIST_NOTAG || xt==XT_LIST_TAG ||
	xt==XT_LANG_NOTAG || xt==XT_LANG_TAG) {
	REXPFactory lc = new REXPFactory();
	REXPFactory nf = new REXPFactory();
	RList l = new RList();
	while (o<eox) {
		String name = null;
		o = lc.parseREXP(buf, o);
		if (xt==XT_LIST_TAG || xt==XT_LANG_TAG) {
			o = nf.parseREXP(buf, o);
			if (nf.cont.isSymbol() || nf.cont.isString()) name = nf.cont.asString();
		}
		if (name==null) l.add(lc.cont);
		else l.put(name, lc.cont);
	}
	cont = (xt==XT_LANG_NOTAG || xt==XT_LANG_TAG)?
		new REXPLanguage(l, getAttr()):
		new REXPList(l, getAttr());
	if (o!=eox) {
		System.err.println("Warning: int list SEXP size mismatch\n");
		o=eox;
	}
	return o;
}

*/
