package assign

import "github.com/senseyeio/roger/constants"

func SetHdr(valueType constants.DataType, valueLength int, buf []byte) {
	buf[0] = byte(valueType)
	buf[1] = byte(valueLength & 255)
	buf[2] = byte((valueLength & 0xff00) >> 8)
	buf[3] = byte((valueLength & 0xff0000) >> 16)
}

func setHdrOffset(valueType constants.DataType, valueLength int, buf []byte, o int) {
	if valueLength > 0xfffff0 {
		buf[o] = byte((valueType & 255) | constants.DtLarge)
		o++
	} else {
		buf[o] = byte(valueType & 255)
		o++
	}
	buf[o] = byte(valueLength & 255)
	o++
	buf[o] = byte((valueLength & 0xff00) >> 8)
	o++
	buf[o] = byte((valueLength & 0xff0000) >> 16)
	o++
	if valueLength > 0xfffff0 {
		buf[o] = byte((valueLength & 0xff000000) >> 24)
		o++
		buf[o] = 0
		o++
		buf[o] = 0
		o++
		buf[o] = 0
		o++
	}
}

func SetInt(v int, buf []byte, o int) {
	buf[o] = byte(v & 255)
	o++
	buf[o] = byte((v & 0xff00) >> 8)
	o++
	buf[o] = byte((v & 0xff0000) >> 16)
	o++
	buf[o] = byte((v & 0xff000000) >> 24)
	o++
}

func setLong(l int64, buf []byte, o int) {
	SetInt(int(l&0xffffffff), buf, o)
	SetInt(int(l>>32), buf, o+4)
}
