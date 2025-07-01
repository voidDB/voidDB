package node

import (
	"github.com/voidDB/voidDB/common"
)

const (
	MaxNodeLength = 7
)

type Node []byte

func NewNode() (node Node) {
	node = common.NewPage()

	node.setMagic()

	return
}

func (node Node) copy() (newNode Node) {
	newNode = common.NewPage()

	copy(newNode, node)

	return newNode
}

func (node Node) meta() Meta {
	return common.LineN(node, 0)
}

func (node Node) setMagic() {
	node.meta().setMagic()

	return
}

func (node Node) VetMagic() error {
	return node.meta().vetMagic()
}

func (node Node) Length() int {
	return node.getLength()
}

func (node Node) getLength() int {
	return node.meta().getLength()
}

func (node Node) setLength(length int) {
	node.meta().setLength(length)

	return
}

func (node Node) elem(index int) Elem {
	return common.LineN(node,
		(index+1)%(MaxNodeLength+1),
	)
}

func (node Node) ValueOrChild(index int) (pointer, length int) {
	return node.getValueOrChild(index)
}

func (node Node) ValueOrChildMeta(index int) []byte {
	return node.elem(index).meta()
}

func (node Node) getValueOrChild(index int) (pointer, length int) {
	var elem Elem = node.elem(index)

	return elem.getPointer(), elem.getValLen()
}

func (node Node) setValueOrChild(index, pointer, length int, elemMeta []byte) {
	var elem Elem = node.elem(index)

	elem.setPointer(pointer)

	elem.setValLen(length)

	elem.setMeta(elemMeta)

	return
}

func (node Node) key(index int) Key {
	return common.KeyN(node, index+1)
}

func (node Node) Key(index int) []byte {
	return node.getKey(index)
}

func (node Node) getKey(index int) []byte {
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

func (node Node) IsGraveyard() bool {
	var (
		index   int
		pointer common.Pointer
	)

	for index = 0; index <= node.getLength(); index++ {
		pointer = common.Pointer(
			node.elem(index).getPointer(),
		)

		switch {
		case pointer.IsDeleted():
			continue

		case pointer == 0 && index == node.getLength():
			continue

		default:
			return false
		}
	}

	return true
}
