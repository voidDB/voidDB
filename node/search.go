package node

import (
	"bytes"
)

func (node Node) Search(key []byte) (index, pointer, length int) {
	var (
		result int
	)

	for index = 0; index < node.Length(); index++ {
		result = bytes.Compare(key,
			node.Key(index),
		)

		if result < 1 {
			break
		}
	}

	pointer, length = node.ValueOrChild(index)

	switch {
	case length > 0 && result == 0: // leaf node, record found
		return

	case length > 0: // leaf node, record not found
		return index, 0, 0
	}

	return
}
