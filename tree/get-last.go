package tree

func (cursor *Cursor) GetLast() (key, value []byte, e error) {
	cursor.reset()

	cursor.index = MaxNodeLength

	return cursor._getPrev()
}
