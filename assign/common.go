package assign

import (
	"github.com/senseyeio/roger/constants"
)

func GetHeaderLength(valueType constants.DataType, valueLength int) int {
	// Simon Urbanek confirmed for pyRserve that this does not cause any problems with Rserve.
	return constants.LgHeaderSize
}

func SetHdr(valueType constants.DataType, valueLength int, buf []byte) {
	setHdrOffset(valueType, valueLength, buf, 0)
}

func setHdrOffset(valueType constants.DataType, valueLength int, buf []byte, o int) {

	// always large headers
	buf[o] = byte((valueType & 255) | constants.DtLarge)
	o++

	// main body
	buf[o] = byte(valueLength & 255)
	o++
	buf[o] = byte((valueLength & 0xff00) >> 8)
	o++
	buf[o] = byte((valueLength & 0xff0000) >> 16)
	o++

	// extra large header content
	buf[o] = byte((int64(valueLength) & 0xff000000) >> 24)
	o++
	buf[o] = 0
	o++
	buf[o] = 0
	o++
	buf[o] = 0
	o++
}

func SetInt(v int, buf []byte, o int) {
	buf[o] = byte(v & 255)
	o++
	buf[o] = byte((v & 0xff00) >> 8)
	o++
	buf[o] = byte((v & 0xff0000) >> 16)
	o++
	buf[o] = byte((int64(v) & 0xff000000) >> 24)
	o++
}

func setLong(l int64, buf []byte, o int) {
	SetInt(int(l&0xffffffff), buf, o)
	SetInt(int(l>>32), buf, o+4)
}
