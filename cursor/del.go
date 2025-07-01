package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

// Del deletes the key-value record indexed by the cursor. To delete a record
// by key, position the cursor with [*Cursor.Get] beforehand.
func (cursor *Cursor) Del() (e error) {
	cursor.resume()

	return cursor.del(
		cursor.medium.Meta(),
	)
}

func (cursor *Cursor) del(metadata []byte) (e error) {
	var (
		newNode node.Node
		oldNode node.Node

		dirty   bool
		i       int
		pointer int
	)

	oldNode, dirty, e = getNode(cursor.medium, cursor.offset, true)
	if e != nil {
		return
	}

	cursor.medium.Free(
		oldNode.ValueOrChild(cursor.index),
	)

	newNode = oldNode.Update(cursor.index, common.Tombstone, -1, metadata,
		dirty,
	)

	pointer = cursor.medium.Save(newNode)

	cursor.offset = pointer

	for i = len(cursor.stack) - 1; i > -1; i-- {
		oldNode, dirty, e = getNode(cursor.medium, cursor.stack[i].offset, true)
		if e != nil {
			return
		}

		if newNode.IsGraveyard() {
			pointer = common.Pointer(pointer).ToGraveyard()
		}

		newNode = oldNode.Update(cursor.stack[i].index, pointer, 0,
			cursor.medium.Meta(),
			dirty,
		)

		pointer = cursor.medium.Save(newNode)

		cursor.stack[i].offset = pointer
	}

	cursor.medium.Root(pointer)

	return
}
