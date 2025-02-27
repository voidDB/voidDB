package cursor

import (
	"math"
	"syscall"

	"github.com/voidDB/voidDB/link"
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
//
// CAUTION: The data in value must not be modified until the transaction has
// been successfully committed.
func (cursor *Cursor) Put(key, value []byte) (e error) {
	return cursor.put(key, value, nil)
}

func (cursor *Cursor) put(key, value, linkMeta link.Metadata) (e error) {
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

	if linkMeta == nil {
		linkMeta = cursor.medium.Meta()
	}

	pointer0, pointer1, promoted, e = put(cursor.medium, cursor.offset,
		cursor.medium.Save(value), len(value), key, linkMeta,
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
			Insert(0, pointer0, pointer1, 0, promoted,
				cursor.medium.Meta(),
			)

		cursor.offset = cursor.medium.Save(newRoot)
	}

	return
}

func put(medium Medium, offset, putPointer, putLength int, key []byte,
	linkMeta link.Metadata,
) (
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

	pointer &^= graveyard

	switch {
	case pointer == 0:
		newNode0, newNode1, promoted = oldNode.Insert(index,
			putPointer, 0, putLength, key, linkMeta,
		)

	case length > 0:
		medium.Free(pointer, length)

		fallthrough

	case pointer == tombstone:
		newNode0 = oldNode.Update(index, putPointer, putLength, linkMeta)

	default:
		pointer0, pointer1, promoted, e = put(medium, pointer,
			putPointer, putLength, key, linkMeta,
		)
		if e != nil {
			return
		}

		switch {
		case pointer1 == 0:
			newNode0 = oldNode.Update(index,
				pointer0, 0, medium.Meta(),
			)

		default:
			newNode0, newNode1, promoted = oldNode.Insert(index,
				pointer0, pointer1, 0, promoted, medium.Meta(),
			)

			pointer1 = 0
		}
	}

	pointer0 = medium.Save(newNode0)

	if newNode1 == nil {
		return
	}

	if isGraveyard(newNode0) {
		pointer0 |= graveyard
	}

	pointer1 = medium.Save(newNode1)

	if isGraveyard(newNode1) {
		pointer1 |= graveyard
	}

	return
}
