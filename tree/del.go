package tree

const (
	tombstone = 1 // Pointers are multiples of PageSize,
	// hence the least significant 11 bits are free to mean other things.
)

func (cursor *Cursor) Del() (e error) {
	cursor.resume()

	return cursor._del()
}

func (cursor *Cursor) _del() (e error) {
	var (
		node Node = getNode(cursor.medium, cursor.offset, true)

		i       int
		pointer int
	)

	cursor.medium.Free(
		node.pointer(cursor.index),
		node.valLen(cursor.index),
	)

	node = node.update(cursor.index, tombstone, 0)

	pointer, e = cursor.medium.Save(node)
	if e != nil {
		return
	}

	cursor.offset = pointer

	for i = len(cursor.stack) - 1; i > -1; i-- {
		node = getNode(cursor.medium, cursor.stack[i].offset, true)

		node = node.update(cursor.stack[i].index, pointer, 0)

		pointer, e = cursor.medium.Save(node)
		if e != nil {
			return
		}

		cursor.stack[i].offset = pointer
	}

	return
}
