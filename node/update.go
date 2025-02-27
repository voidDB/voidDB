package node

import (
	"github.com/voidDB/voidDB/link"
)

func (node Node) Update(index, pointer, length int, metadata link.Metadata) (
	newNode Node,
) {
	newNode = make([]byte, pageSize)

	copy(newNode, node)

	switch {
	case length < 0:
		newNode.elem(index).setPointer(pointer)

		newNode.elem(index).setLinkMetadata(metadata)

	default:
		newNode.setValueOrChild(index, pointer, length, metadata)
	}

	return
}
