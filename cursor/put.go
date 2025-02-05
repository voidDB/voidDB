package cursor

import (
	"math"
	"syscall"

	"github.com/voidDB/voidDB/node"
)

const (
	MaxKeyLength   = node.MaxKeyLength
	MaxValueLength = math.MaxUint32
)

// Put stores a key-value pair (or overwrites the corresponding value, if key
// already exists) and positions the cursor at the inserted record. It returns
// [syscall.EINVAL] (“invalid argument”) if the length of key or value is zero,
// or otherwise exceeds [MaxKeyLength] or [MaxValueLength] respectively.
func (cursor *Cursor) Put(key, value []byte) (e error) {
	var (
		newRoot  node.Node
		pointer0 int
		pointer1 int
		promoted []byte
	)

	switch {
	case len(key) == 0:
		fallthrough

	case len(key) > MaxKeyLength:
		fallthrough

	case len(value) == 0:
		fallthrough

	case len(value) > MaxValueLength:
		return syscall.EINVAL
	}

	cursor.reset()

	pointer0, pointer1, promoted, e = put(cursor.medium, cursor.offset,
		cursor.medium.Save(value),
		len(value),
		key,
	)
	if e != nil {
		return
	}

	cursor.latest = key

	switch {
	case pointer1 == 0:
		cursor.offset = pointer0

	default:
		newRoot, _, _ = node.NewNode().
			Insert(0, pointer0, pointer1, 0, promoted)

		cursor.offset = cursor.medium.Save(newRoot)
	}

	return
}

func put(medium Medium, offset, putPointer, putLength int, key []byte) (
	pointer0, pointer1 int, promoted []byte, e error,
) {
	var (
		oldNode  node.Node
		newNode0 node.Node
		newNode1 node.Node

		index   int
		length  int
		pointer int
	)

	oldNode, e = getNode(medium, offset, true)
	if e != nil {
		return
	}

	index, pointer, length = oldNode.Search(key)

	switch {
	case pointer == 0:
		newNode0, newNode1, promoted = oldNode.Insert(index,
			putPointer,
			0,
			putLength,
			key,
		)

	case length > 0:
		medium.Free(pointer, length)

		fallthrough

	case pointer == tombstone:
		newNode0 = oldNode.Update(index, putPointer, putLength)

	default:
		pointer0, pointer1, promoted, e = put(medium, pointer,
			putPointer,
			putLength,
			key,
		)
		if e != nil {
			return
		}

		switch {
		case pointer1 == 0:
			newNode0 = oldNode.Update(index, pointer0, 0)

		default:
			newNode0, newNode1, promoted = oldNode.Insert(index,
				pointer0,
				pointer1,
				0,
				promoted,
			)

			pointer1 = 0
		}
	}

	pointer0 = medium.Save(newNode0)

	if newNode1 != nil {
		pointer1 = medium.Save(newNode1)
	}

	return
}
