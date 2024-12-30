package tree

func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	var (
		node Node = getNode(cursor.medium, cursor.offset, false)

		pointer int = node.pointer(cursor.index)
		valLen  int = node.valLen(cursor.index)
	)

	switch {
	case pointer == 0 && len(cursor.stack) == 0:
		e = ErrorNotFound

		return

	case pointer == 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index+1

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.GetNext()

	case valLen > 0:
		defer func() { cursor.index++ }()

		return node.key(cursor.index),
			cursor.medium.Load(pointer, valLen),
			nil
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index},
	)

	cursor.offset, cursor.index = pointer, 0

	return cursor.GetNext()
}

type ancestor struct {
	offset int
	index  int
}
