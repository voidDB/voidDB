package cursor

import (
	"errors"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

// GetPrev regresses the cursor and retrieves the previous key-value record,
// sorted by key using [bytes.Compare].
//
// CAUTION: See [*Cursor.Get].
func (cursor *Cursor) GetPrev() (key, value []byte, e error) {
	cursor.resume()

	for {
		key, value, e = cursor.getPrev()

		if !errors.Is(e, common.ErrorDeleted) {
			break
		}
	}

	return
}

func (cursor *Cursor) getPrev() (key, value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	cursor.index--

	switch {
	case cursor.index < 0 && len(cursor.stack) == 0:
		return nil, nil, common.ErrorNotFound

	case cursor.index < 0:
		cursor.ascend()

		return cursor.getPrev()
	}

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	if cursor.index == node.MaxNodeLength {
		cursor.index = curNode.Length()
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	switch {
	case common.Pointer(pointer).IsDeleted():
		return nil, nil, common.ErrorDeleted

	case length > 0:
		key = curNode.Key(cursor.index)

		value = cursor.medium.Data(pointer, length)

		return

	case pointer > 0:
		cursor.descend(pointer, node.MaxNodeLength)
	}

	return cursor.getPrev()
}
