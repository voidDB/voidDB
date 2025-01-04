package tree

func (cursor *Cursor) GetPrev() (key, value []byte, e error) {
	cursor.index--

	return cursor._getPrev()
}

func (cursor *Cursor) _getPrev() (key, value []byte, e error) {
	var (
		node    Node
		pointer int
		valLen  int
	)

	switch {
	case cursor.index < 0 && len(cursor.stack) == 0:
		e = ErrorNotFound

		return

	case cursor.index < 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index-1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor._getPrev()
	}

	node = getNode(cursor.medium, cursor.offset, false)

	if cursor.index == MaxNodeLength {
		cursor.index = node.length() - 1

		return cursor._getPrev()
	}

	pointer, valLen = node.pointer(cursor.index), node.valLen(cursor.index)

	switch {
	case valLen > 0:
		return node.key(cursor.index),
			cursor.medium.Load(pointer, valLen),
			nil

	case pointer == tombstone:
		return cursor.GetPrev()
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset, cursor.index = pointer, MaxNodeLength

	return cursor._getPrev()
}
