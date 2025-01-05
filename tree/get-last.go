package tree

func (cursor *Cursor) GetLast() (key, value []byte, e error) {
	cursor.reset()

	cursor.index = maxNodeLength

	return cursor._getPrev()
}
