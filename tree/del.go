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
		i       int
		node    Node
		pointer int
	)

	node, e = getNode(cursor.medium, cursor.offset, true)
	if e != nil {
		return
	}

	cursor.medium.Free(
		node.pointer(cursor.index),
		node.valLen(cursor.index),
	)

	node = node.update(cursor.index, tombstone, 0)

	pointer = cursor.medium.Save(node)

	cursor.offset = pointer

	for i = len(cursor.stack) - 1; i > -1; i-- {
		node, e = getNode(cursor.medium, cursor.stack[i].offset, true)
		if e != nil {
			return
		}

		node = node.update(cursor.stack[i].index, pointer, 0)

		pointer = cursor.medium.Save(node)

		cursor.stack[i].offset = pointer
	}

	return
}
