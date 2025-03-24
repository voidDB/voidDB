package cursor

import (
	"github.com/voidDB/voidDB/link"
	"github.com/voidDB/voidDB/node"
)

const (
	tombstone = 1 // HACK: Pointers are multiples of common.WordSize, hence the
	graveyard = 2 // least significant 3 bits are free to mean other things.
)

// Del deletes the key-value record indexed by the cursor. To delete a record
// by key, position the cursor with [*Cursor.Get] beforehand.
func (cursor *Cursor) Del() (e error) {
	cursor.resume()

	return cursor.del(nil)
}

func (cursor *Cursor) del(linkMeta link.Metadata) (e error) {
	var (
		newNode node.Node
		oldNode node.Node

		dirty   bool
		i       int
		pointer int
	)

	if linkMeta == nil {
		linkMeta = cursor.medium.Meta()
	}

	oldNode, dirty, e = getNode(cursor.medium, cursor.offset, true)
	if e != nil {
		return
	}

	cursor.medium.Free(
		oldNode.ValueOrChild(cursor.index),
	)

	newNode = oldNode.Update(cursor.index, tombstone, -1, linkMeta, dirty)

	pointer = cursor.medium.Save(newNode)

	cursor.offset = pointer

	for i = len(cursor.stack) - 1; i > -1; i-- {
		oldNode, dirty, e = getNode(cursor.medium, cursor.stack[i].offset, true)
		if e != nil {
			return
		}

		if isGraveyard(newNode) {
			pointer |= graveyard
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

func isGraveyard(n node.Node) bool {
	var (
		index   int
		pointer int
	)

	for index = 0; index <= n.Length(); index++ {
		pointer, _ = n.ValueOrChild(index)

		switch {
		case pointer&graveyard > 0:
			continue

		case pointer == tombstone:
			continue

		case pointer == 0 && index == n.Length():
			continue

		default:
			return false
		}
	}

	return true
}
