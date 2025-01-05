package tree

func (cursor *Cursor) GetPrev() (key, value []byte, e error) {
	cursor.resume()

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

	if cursor.index == maxNodeLength {
		cursor.index = node.length()
	}

	pointer = node.pointer(cursor.index)

	if pointer == 0 || pointer == tombstone {
		return cursor.GetPrev()
	}

	valLen = node.valLen(cursor.index)

	if valLen > 0 {
		return node.key(cursor.index),
			cursor.medium.Load(pointer, valLen),
			nil
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset, cursor.index = pointer, maxNodeLength

	return cursor._getPrev()
}
