package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	cursor.resume()

	cursor.index++

	return cursor.getNext()
}

func (cursor *Cursor) getNext() (key, value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	curNode, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	switch {
	case length > 0:
		key = curNode.Key(cursor.index)

		value = cursor.medium.Load(pointer, length)

		return

	case pointer == tombstone:
		cursor.index++

		return cursor.getNext()

	case pointer > 0:
		cursor.stack = append(cursor.stack,
			ancestor{cursor.offset, cursor.index},
		)

		cursor.offset, cursor.index = pointer, 0

		return cursor.getNext()

	case len(cursor.stack) > 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index+1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.getNext()
	}

	e = common.ErrorNotFound

	return
}
