package sexp

import (
	"errors"
	"log"
)

func parseVectorAttr(attr interface{}, vectorArr []interface{}, offset int) (interface{}, int, error) {
	attrMap, ok := attr.(map[string]interface{})
	if !ok {
		return vectorArr, offset, nil
	}
	names, ok := attrMap["n"].([]string)
	if !ok {
		return vectorArr, offset, nil
	}
	if len(names) != len(vectorArr) {
		return nil, offset, errors.New("Vector name and value quantity mismatch")
	}

	ret := map[string]interface{}{}
	for i := range names {
		ret[names[i]] = vectorArr[i]
	}

	return ret, offset, nil
}

func parseVector(attr interface{}, buf []byte, offset, end int) (interface{}, int, error) {
	var vectorArr []interface{}

	for offset < end {
		vectorEntry, newOffset, err := parseReturningOffset(buf, offset)
		offset = newOffset
		if err != nil {
			vectorArr = append(vectorArr, err)
		} else {
			vectorArr = append(vectorArr, vectorEntry)
		}
	}
	if offset != end {
		log.Println("Warning: vector size mismatch")
		offset = end
	}
	return parseVectorAttr(attr, vectorArr, offset)
}

/*
if (xt==XT_VECTOR || xt==XT_VECTOR_EXP) {
    Vector v=new Vector();
    while(o<eox) {
        REXPFactory xx=new REXPFactory();
        o = xx.parseREXP(buf,o);
        v.addElement(xx.cont);
    }
    if (o!=eox) {
        System.err.println("Warning: int vector SEXP size mismatch\n");
        o=eox;
    }
    // fixup for lists since they're stored as attributes of vectors
    if (getAttr()!=null && getAttr().asList().at("names") != null) {
        REXP nam = getAttr().asList().at("names");
        String names[] = null;
        if (nam.isString()) names = nam.asStrings();
        else if (nam.isVector()) { // names could be a vector if supplied by old Rserve
            RList l = nam.asList();
            Object oa[] = l.toArray();
            names = new String[oa.length];
            for(int i = 0; i < oa.length; i++) names[i] = ((REXP)oa[i]).asString();
        }
        RList l = new RList(v, names);
        cont = (xt==XT_VECTOR_EXP)?
            new REXPExpressionVector(l, getAttr()):
            new REXPGenericVector(l, getAttr());
    } else
        cont = (xt==XT_VECTOR_EXP)?
            new REXPExpressionVector(new RList(v), getAttr()):
            new REXPGenericVector(new RList(v), getAttr());
    return o;
}

*/
