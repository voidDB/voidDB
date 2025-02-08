package node

import (
	"github.com/voidDB/voidDB/common"
)

const (
	MaxKeyLength  = 512
	MaxNodeLength = 7

	pageSize = common.PageSize
	lineSize = common.LineSize
	wordSize = common.WordSize
	halfSize = common.HalfSize

	nodeMagic = "voidNODE"
)

type Node []byte

func NewNode() (node Node) {
	node = make([]byte, pageSize)

	node.meta().setMagic()

	return
}

func (node Node) meta() Meta {
	return common.Field(node, 0, lineSize)
}

func (node Node) IsNode() bool {
	return node.meta().isNode()
}

func (node Node) Length() int {
	return node.meta().getLength()
}

func (node Node) setLength(length int) {
	node.meta().setLength(length)

	return
}

func (node Node) elem(index int) Elem {
	return common.Field(node,
		lineSize*((index+1)%(MaxNodeLength+1)),
		lineSize,
	)
}

func (node Node) ValueOrChild(index int) (pointer, length int) {
	var elem Elem = node.elem(index)

	return elem.getPointer(), elem.getValLen()
}

func (node Node) setValueOrChild(index, pointer, length int, metadata []byte) {
	var elem Elem = node.elem(index)

	elem.setPointer(pointer)

	elem.setValLen(length)

	elem.setExtraMetadata(metadata)

	return
}

func (node Node) key(index int) Key {
	return common.Field(node,
		MaxKeyLength*(index+1),
		MaxKeyLength,
	)
}

func (node Node) Key(index int) []byte {
	return node.key(index).get(
		node.elem(index).getKeyLen(),
	)
}

func (node Node) setKey(index int, key []byte) {
	node.elem(index).setKeyLen(
		node.key(index).set(key),
	)

	return
}
