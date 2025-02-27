package cursor

import (
	"github.com/voidDB/voidDB/link"
)

type linkCursor struct {
	*Cursor
}

func (cursor *Cursor) ToLinkCursor() *linkCursor {
	return &linkCursor{cursor}
}

func (cursor *linkCursor) GetNext(minTxnID int) (
	key, value []byte, linkMeta link.Metadata, e error,
) {
	return cursor.getNextWithLeafMeta(minTxnID)
}

func (cursor *linkCursor) Get(key []byte) (
	linkMeta link.Metadata, e error,
) {
	return cursor.getLeafMetaReset(key)
}

func (cursor *linkCursor) Put(key, value []byte, linkMeta link.Metadata) (
	e error,
) {
	return cursor.put(key, value, linkMeta)
}

func (cursor *linkCursor) Del(linkMeta link.Metadata) (e error) {
	return cursor.del(linkMeta)
}
