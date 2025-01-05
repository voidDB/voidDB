package tree

func (cursor *Cursor) Get(key []byte) (value []byte, e error) {
	cursor.reset()

	return cursor._get(key)
}

func (cursor *Cursor) _get(key []byte) (value []byte, e error) {
	var (
		node Node = getNode(cursor.medium, cursor.offset, false)

		pointer int
		valLen  int
	)

	cursor.index, pointer, valLen = node.search(key)

	switch {
	case pointer == tombstone:
		fallthrough

	case pointer == 0:
		return nil, ErrorNotFound

	case valLen > 0:
		return cursor.medium.Load(pointer, valLen), nil
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset = pointer

	return cursor._get(key)
}
