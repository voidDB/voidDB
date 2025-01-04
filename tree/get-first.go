package tree

func (cursor *Cursor) GetFirst() (key, value []byte, e error) {
	cursor.reset()

	return cursor.GetNext()
}
