package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func (cursor *Cursor) GetPrev() (key, value []byte, e error) {
	cursor.resume()

	cursor.index--

	return cursor.getPrev()
}

func (cursor *Cursor) getPrev() (key, value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	switch {
	case cursor.index < 0 && len(cursor.stack) == 0:
		e = common.ErrorNotFound

		return

	case cursor.index < 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index-1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.getPrev()
	}

	curNode, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	if cursor.index == node.MaxNodeLength {
		cursor.index = curNode.Length()
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	switch {
	case pointer == 0 || pointer == tombstone:
		cursor.index--

		return cursor.getPrev()

	case length > 0:
		key = curNode.Key(cursor.index)

		value = cursor.medium.Load(pointer, length)

		return
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset, cursor.index = pointer, node.MaxNodeLength

	return cursor.getPrev()
}
