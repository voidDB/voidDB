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

	cursor.index++

	if cursor.index >= node.MaxNodeLength {
		goto end
	}

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
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
		cursor.descend(pointer, -1)

		return cursor.getNext()
	}

end:
	if len(cursor.stack) > 0 {
		cursor.ascend()

		return cursor.getNext()
	}

	return nil, nil, common.ErrorNotFound
}

func (cursor *Cursor) getNextWithLeafMeta(minTxnID int) (
	key, value []byte, meta link.Metadata, e error,
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

	meta = curNode.ValueOrChildMeta(cursor.index)

	if meta.TxnSerial().Int() < minTxnID {
		return cursor.getNextWithLeafMeta(minTxnID)
	}

	pointer, length = curNode.ValueOrChild(cursor.index)

	switch {
	case common.Pointer(pointer).IsTombstone():
		e = common.ErrorDeleted

		fallthrough

	case length > 0:
		key = curNode.Key(cursor.index)

		if e == nil {
			value = cursor.medium.Data(pointer, length)
		}

		return

	case pointer > 0:
		cursor.descend(
			common.Pointer(pointer).Clean(),
			-1,
		)

		return cursor.getNextWithLeafMeta(minTxnID)
	}

end:
	if len(cursor.stack) > 0 {
		cursor.ascend()

		return cursor.getNextWithLeafMeta(minTxnID)
	}

	return nil, nil, nil, common.ErrorNotFound
}
