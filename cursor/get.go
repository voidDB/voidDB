package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/link"
	"github.com/voidDB/voidDB/node"
)

// Get retrieves the value corresponding to key and positions the cursor at the
// record found.
//
// CAUTION: The value returned is immutable and valid only during the lifetime
// of the transaction to which the cursor belongs, since the slice merely
// reflects the relevant section of the memory map containing the value. Hence,
// any attempt at mutating the slice at any time or accessing it after the
// transaction has been committed/aborted will result in a fatal
// [syscall.SIGSEGV]. (See also [runtime/debug.SetPanicOnFault].) Instead,
// applications should allocate a slice of size equal to len(value) and copy
// value into the new slice for modification/retention. This also applies to
// [*Cursor.GetFirst], [*Cursor.GetLast], [*Cursor.GetNext], and
// [*Cursor.GetPrev].
func (cursor *Cursor) Get(key []byte) (value []byte, e error) {
	cursor.reset()

	return cursor.get(key)
}

// GetFirst retrieves and positions the cursor at the first key-value record,
// sorted by key using [bytes.Compare].
//
// CAUTION: See [*Cursor.Get].
func (cursor *Cursor) GetFirst() (key, value []byte, e error) {
	cursor.reset()

	return cursor.GetNext()
}

// GetLast retrieves and positions the cursor at the last key-value record,
// sorted by key using [bytes.Compare].
//
// CAUTION: See [*Cursor.Get].
func (cursor *Cursor) GetLast() (key, value []byte, e error) {
	cursor.reset()

	cursor.index = node.MaxNodeLength

	return cursor.GetPrev()
}

func (cursor *Cursor) get(key []byte) (value []byte, e error) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	cursor.index, pointer, length = curNode.Search(key)

	switch {
	case pointer&graveyard > 0:
		fallthrough

	case pointer == tombstone:
		fallthrough

	case pointer == 0:
		return nil, common.ErrorNotFound

	case length > 0:
		value, _ = cursor.medium.Load(pointer, length)

		return
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset = pointer

	return cursor.get(key)
}

func (cursor *Cursor) getLeafMetaReset(key []byte) (
	linkMeta link.Metadata, e error,
) {
	cursor.reset()

	return cursor.getLeafMeta(key)
}

func (cursor *Cursor) getLeafMeta(key []byte) (
	linkMeta link.Metadata, e error,
) {
	var (
		curNode node.Node
		length  int
		pointer int
	)

	curNode, _, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	cursor.index, pointer, length = curNode.Search(key)

	switch {
	case pointer&graveyard > 0:
		break

	case pointer == tombstone:
		e = common.ErrorDeleted

		fallthrough

	case length > 0:
		linkMeta = curNode.ValueOrChildLinkMetadata(cursor.index)

		return

	case pointer == 0:
		return nil, common.ErrorNotFound
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset = pointer &^ graveyard

	return cursor.getLeafMeta(key)
}
