package node

import (
	"github.com/voidDB/voidDB/common"
)

type Elem []byte

func (elem Elem) keyLen() []byte {
	return common.HalfN(elem, 0)
}

func (elem Elem) getKeyLen() int {
	return common.GetIntFromHalf(
		elem.keyLen(),
	)
}

func (elem Elem) setKeyLen(length int) {
	common.PutIntIntoHalf(
		elem.keyLen(),
		length,
	)

	return
}

func (elem Elem) valLen() []byte {
	return common.HalfN(elem, 1)
}

func (elem Elem) getValLen() int {
	return common.GetIntFromHalf(
		elem.valLen(),
	)
}

func (elem Elem) setValLen(length int) {
	common.PutIntIntoHalf(
		elem.valLen(),
		length,
	)

	return
}

func (elem Elem) pointer() []byte {
	return common.WordN(elem, 1)
}

func (elem Elem) getPointer() int {
	return common.GetIntFromWord(
		elem.pointer(),
	)
}

func (elem Elem) setPointer(pointer int) {
	common.PutIntIntoWord(
		elem.pointer(),
		pointer,
	)

	return
}

func (elem Elem) meta() []byte {
	return common.TwinN(elem, 1)
}

func (elem Elem) getMeta() []byte {
	return elem.meta()
}

func (elem Elem) setMeta(meta []byte) {
	copy(
		elem.meta(),
		meta,
	)

	return
}
