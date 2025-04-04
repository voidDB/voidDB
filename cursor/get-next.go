package cursor

import (
	"errors"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/link"
	"github.com/voidDB/voidDB/node"
)

// GetNext advances the cursor and retrieves the next key-value record, sorted
// by key using [bytes.Compare]. On a newly opened cursor, it has the same
// effect as [*Cursor.GetFirst].
//
// CAUTION: See [*Cursor.Get].
func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	cursor.resume()

	for {
		cursor.index++

		key, value, e = cursor.getNext()

		if !errors.Is(e, common.ErrorDeleted) {
			break
		}
	}

	return
}

func (cursor *Cursor) getNext() (key, value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	if cursor.index >= node.MaxNodeLength {
		goto end
	}

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	switch {
	case pointer&graveyard > 0:
		fallthrough

	case pointer == tombstone:
		return nil, nil, common.ErrorDeleted

	case length > 0:
		key = curNode.Key(cursor.index)

		value, _ = cursor.medium.Load(pointer, length)

		return

	case pointer > 0:
		cursor.stack = append(cursor.stack,
			ancestor{cursor.offset, cursor.index},
		)

		cursor.offset, cursor.index = pointer, 0

		return cursor.getNext()
	}

end:
	if len(cursor.stack) > 0 {
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index+1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.getNext()
	}

	e = common.ErrorNotFound

	return
}

func (cursor *Cursor) getNextWithLeafMeta(minTxnID int) (
	key, value []byte, linkMeta link.Metadata, e error,
) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	cursor.index++

	if cursor.index >= node.MaxNodeLength {
		goto end
	}

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	linkMeta = curNode.ValueOrChildLinkMetadata(cursor.index)

	if linkMeta.TxnSerial().Int() < minTxnID {
		return cursor.getNextWithLeafMeta(minTxnID)
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	pointer &^= graveyard

	switch {
	case pointer == tombstone:
		e = common.ErrorDeleted

		fallthrough

	case length > 0:
		key = curNode.Key(cursor.index)

		if e == nil {
			value, _ = cursor.medium.Load(pointer, length)
		}

		return

	case pointer > 0:
		cursor.stack = append(cursor.stack,
			ancestor{cursor.offset, cursor.index},
		)

		cursor.offset, cursor.index = pointer, 0

		return cursor.getNextWithLeafMeta(minTxnID)
	}

end:
	if len(cursor.stack) > 0 {
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index+1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.getNextWithLeafMeta(minTxnID)
	}

	e = common.ErrorNotFound

	return
}
