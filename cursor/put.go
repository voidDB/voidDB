package cursor

import (
	"math"
	"syscall"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

const (
	MaxKeyLength   = common.KeySize
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
	return cursor.put(key, value,
		cursor.medium.Meta(),
	)
}

func (cursor *Cursor) put(key, value, metadata []byte) (e error) {
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
		cursor.medium.Save(value), len(value), key, metadata,
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
				cursor.medium.Meta(), true,
			)

		cursor.offset = cursor.medium.Save(newRoot)
	}

	cursor.medium.Root(cursor.offset)

	return
}

func put(medium Medium, offset, putPointer, putLength int, key, metadata []byte,
) (
	pointer0, pointer1 int, promoted []byte, e error,
) {
	var (
		newNode0 node.Node
		newNode1 node.Node
		oldNode  node.Node

		dirty   bool
		index   int
		length  int
		pointer int
		replace bool
	)

	oldNode, dirty, e = getNode(medium, offset, true)
	if e != nil {
		return
	}

	index, pointer, length = oldNode.Search(key)

	switch {
	case common.Pointer(pointer).IsTombstone():
		replace = true

	case length > 0:
		medium.Free(pointer, length)

		replace = true

	case pointer > 0:
		pointer0, pointer1, promoted, e = put(medium,
			common.Pointer(pointer).Clean(),
			putPointer, putLength, key, metadata,
		)
		if e != nil {
			return
		}

		switch {
		case pointer1 == 0:
			newNode0 = oldNode.Update(index, pointer0, 0, medium.Meta(), dirty)

		default:
			newNode0, newNode1, promoted = oldNode.Insert(index,
				pointer0, pointer1, 0, promoted, medium.Meta(), dirty,
			)

			pointer1 = 0
		}

	default:
		newNode0, newNode1, promoted = oldNode.Insert(index,
			putPointer, 0, putLength, key, metadata, dirty,
		)
	}

	if replace {
		newNode0 = oldNode.Update(index, putPointer, putLength, metadata, dirty)
	}

	pointer0 = medium.Save(newNode0)

	if newNode1 == nil {
		return
	}

	if newNode0.IsGraveyard() {
		pointer0 = common.Pointer(pointer0).ToGraveyard()
	}

	pointer1 = medium.Save(newNode1)

	if newNode1.IsGraveyard() {
		pointer1 = common.Pointer(pointer1).ToGraveyard()
	}

	return
}
