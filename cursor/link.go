package cursor

type linkCursor struct {
	*Cursor
}

func (cursor *Cursor) ToLinkCursor() *linkCursor {
	return &linkCursor{cursor}
}

func (cursor *linkCursor) GetNext() (key, value, leafMeta []byte, e error) {
	return cursor.getNextWithLeafMeta()
}

func (cursor *linkCursor) Get(key []byte) (leafMeta []byte, e error) {
	return cursor.getLeafMetaReset(key)
}

func (cursor *linkCursor) Put(key, value, leafMeta []byte) (e error) {
	return cursor.put(key, value, leafMeta)
}

func (cursor *linkCursor) Del(leafMeta []byte) (e error) {
	return cursor.del(leafMeta)
}
