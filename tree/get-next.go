package tree

func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	cursor.resume()

	cursor.index++

	return cursor._getNext()
}

func (cursor *Cursor) _getNext() (key, value []byte, e error) {
	var (
		node    Node
		pointer int
		valLen  int
	)

	node, e = getNode(cursor.medium, cursor.offset, false)
	if e != nil {
		return
	}

	pointer, valLen = node.pointer(cursor.index), node.valLen(cursor.index)

	switch {
	case valLen > 0:
		key, value = node.key(cursor.index), cursor.medium.Load(pointer, valLen)

		return

	case pointer == tombstone:
		return cursor.GetNext()

	case pointer > 0:
		cursor.stack = append(cursor.stack,
			ancestor{cursor.offset, cursor.index},
		)

		cursor.offset, cursor.index = pointer, 0

		return cursor._getNext()

	case len(cursor.stack) > 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index+1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor._getNext()
	}

	e = errorNotFound

	return
}
