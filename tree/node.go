package tree

import (
	"encoding/binary"
)

const (
	lengthSize  = 4
	magicSize   = 8
	pointerSize = 8

	nodeMagic = "voidNODE"
)

type Node []byte

func (node *Node) _magic() []byte {
	return (*node)[:magicSize]
}

func (node *Node) setMagic() {
	copy(
		node._magic(),
		[]byte(nodeMagic),
	)

	return
}

func (node *Node) _keyLen(index int) []byte {
	const (
		offset = 8
	)

	return (*node)[offset+lengthSize*index : offset+lengthSize*(index+1)]
}

func (node *Node) keyLen(index int) int {
	return int(
		binary.BigEndian.Uint32(
			node._keyLen(index),
		),
	)
}

func (node *Node) setKeyLen(index, keyLen int) {
	binary.BigEndian.PutUint32(
		node._keyLen(index),
		uint32(keyLen),
	)

	return
}

func (node *Node) _valLen(index int) []byte {
	const (
		offset = 36
	)

	return (*node)[offset+lengthSize*index : offset+lengthSize*(index+1)]
}

func (node *Node) valLen(index int) int {
	return int(
		binary.BigEndian.Uint32(
			node._valLen(index),
		),
	)
}

func (node *Node) setValLen(index, valLen int) {
	binary.BigEndian.PutUint32(
		node._valLen(index),
		uint32(valLen),
	)

	return
}

func (node *Node) _pointer(index int) []byte {
	const (
		offset = 64
	)

	return (*node)[offset+pointerSize*index : offset+pointerSize*(index+1)]
}

func (node *Node) pointer(index int) int {
	return int(
		binary.BigEndian.Uint64(
			node._pointer(index),
		),
	)
}

func (node *Node) setPointer(index, pointer int) {
	binary.BigEndian.PutUint64(
		node._pointer(index),
		uint64(pointer),
	)

	return
}

func (node *Node) _key(index int) []byte {
	const (
		offset = 512
	)

	return (*node)[offset+MaxKeySize*index : offset+MaxKeySize*(index+1)]
}

func (node *Node) key(index int) []byte {
	return node._key(index)[:node.keyLen(index)]
}

func (node *Node) setKey(index int, key []byte) {
	var (
		keyLen int = copy(
			node._key(index),
			key,
		)
	)

	node.setKeyLen(index, keyLen)

	return
}

func (node *Node) length() (l int) {
	for l = 0; l < MaxNodeLength; l++ {
		if node.keyLen(l) == 0 {
			break
		}
	}

	return
}
