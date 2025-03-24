package node

import (
	"github.com/voidDB/voidDB/link"
)

func (node Node) Insert(
	index, pointerL, pointerR, length int, key, metadata link.Metadata,
	inPlace bool,
) (
	newNode, _ Node, _ []byte,
) {
	var (
		i          int
		nodeLength int = node.Length()
	)

	switch {
	case inPlace:
		newNode = node

	default:
		newNode = NewNode()

		for i = 0; i < index; i++ {
			copyElemKey(newNode, node, i, 0)
		}
	}

	copyElem(newNode, node, nodeLength, 1)

	for i = nodeLength - 1; i >= index; i-- {
		copyElemKey(newNode, node, i, 1)
	}

	newNode.setKey(index, key)

	newNode.setValueOrChild(index, pointerL, length, metadata)

	if pointerR > 0 {
		newNode.setValueOrChild(index+1, pointerR, length, metadata)
	}

	if nodeLength+1 == MaxNodeLength {
		return newNode.split()
	}

	newNode.setLength(
		nodeLength + 1,
	)

	return
}

func (node Node) split() (newNodeL, newNodeR Node, promoted []byte) {
	var (
		i     int
		shift int = -MaxNodeLength / 2
	)

	newNodeL = NewNode()
	newNodeR = NewNode()

	for i = 0; i < MaxNodeLength; i++ {
		switch {
		case i < MaxNodeLength/2:
			copyElemKey(newNodeL, node, i, 0)

		case i == MaxNodeLength/2 && node.elem(i).getValLen() == 0:
			copyElem(newNodeL, node, i, 0)

			promoted = node.Key(i)

			shift -= 1

		case i == MaxNodeLength/2:
			promoted = node.Key(i - 1)

			fallthrough

		default:
			copyElemKey(newNodeR, node, i, shift)
		}
	}

	copyElem(newNodeR, node, i, shift)

	newNodeL.setLength(MaxNodeLength / 2)

	newNodeR.setLength(MaxNodeLength + shift)

	return
}
