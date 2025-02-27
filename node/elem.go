package node

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/link"
)

type Elem []byte

func (elem Elem) keyLen() []byte {
	return common.Field(elem, 0, halfSize)
}

func (elem Elem) getKeyLen() int {
	return common.GetInt(
		elem.keyLen(),
	)
}

func (elem Elem) setKeyLen(length int) {
	common.PutInt(
		elem.keyLen(),
		length,
	)

	return
}

func (elem Elem) valLen() []byte {
	return common.Field(elem, halfSize, halfSize)
}

func (elem Elem) getValLen() int {
	return common.GetInt(
		elem.valLen(),
	)
}

func (elem Elem) setValLen(length int) {
	common.PutInt(
		elem.valLen(),
		length,
	)

	return
}

func (elem Elem) pointer() []byte {
	return common.Field(elem, wordSize, wordSize)
}

func (elem Elem) getPointer() int {
	return common.GetInt(
		elem.pointer(),
	)
}

func (elem Elem) setPointer(pointer int) {
	common.PutInt(
		elem.pointer(),
		pointer,
	)

	return
}

func (elem Elem) linkMetadata() link.Metadata {
	return common.Field(elem, 2*wordSize, 2*wordSize)
}

func (elem Elem) setLinkMetadata(metadata link.Metadata) {
	copy(
		elem.linkMetadata(),
		metadata,
	)

	return
}
