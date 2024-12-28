package tree

import (
	"bytes"
)

func (node *Node) search(key []byte) (index, pointer, valLen int) {
	var (
		result int
	)

	index, result = node._search(key)

	valLen = node.valLen(index)

	switch {
	case valLen > 0 && result == 0: // leaf node, record found
		break

	case valLen > 0: // leaf node, record not found
		valLen = 0

		return
	}

	pointer = node.pointer(index)

	return
}

func (node *Node) _search(key []byte) (index, result int) {
	for index = 0; index < node.length(); index++ {
		result = bytes.Compare(key,
			node.key(index),
		)

		if result > 0 {
			continue
		}

		break
	}

	return
}
