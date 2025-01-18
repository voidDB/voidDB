package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func (cursor *Cursor) Get(key []byte) (value []byte, e error) {
	cursor.reset()

	return cursor.get(key)
}

func (cursor *Cursor) GetFirst() (key, value []byte, e error) {
	cursor.reset()

	return cursor.GetNext()
}

func (cursor *Cursor) GetLast() (key, value []byte, e error) {
	cursor.reset()

	cursor.index = node.MaxNodeLength

	return cursor.getPrev()
}

func (cursor *Cursor) get(key []byte) (value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	curNode, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	cursor.index, pointer, length = curNode.Search(key)

	switch {
	case pointer == tombstone:
		fallthrough

	case pointer == 0:
		return nil, common.ErrorNotFound

	case length > 0:
		return cursor.medium.Load(pointer, length), nil
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset = pointer

	return cursor.get(key)
}
