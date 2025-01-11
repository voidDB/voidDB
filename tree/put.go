package tree

func put(medium Medium, offset int, key, value []byte) (pointer int, e error) {
	var (
		newRoot  Node
		pointer1 int
		promoted []byte
	)

	pointer = medium.Save(value)

	pointer, pointer1, promoted, e = _put(
		medium, offset, key, pointer, len(value),
	)
	if e != nil {
		return
	}

	if pointer1 == 0 {
		return
	}

	newRoot = NewNode()

	newRoot = newRoot.insert(0, pointer, 0, promoted)

	newRoot.setPointer(1, pointer1)

	pointer = medium.Save(newRoot)

	return
}

func _put(medium Medium, offset int, key []byte, putPointer, putValLen int) (
	pointer, pointer1 int, promoted []byte, e error,
) {
	var (
		index  int
		node   Node
		node1  Node
		node2  Node
		valLen int
	)

	node, e = getNode(medium, offset, true)
	if e != nil {
		return
	}

	index, pointer, valLen = node.search(key)

	switch {
	case pointer == 0:
		node = node.insert(index, putPointer, putValLen, key)

	case valLen > 0:
		medium.Free(pointer, valLen)

		fallthrough

	case pointer == tombstone:
		node = node.update(index, putPointer, putValLen)

	default:
		pointer, pointer1, promoted, e = _put(
			medium, pointer, key, putPointer, putValLen,
		)
		if e != nil {
			return
		}

		switch {
		case pointer1 == 0:
			node = node.update(index, pointer, 0)

		default:
			node = node.insert(index, pointer, 0, promoted)

			node.setPointer(index+1, pointer1)

			pointer1 = 0
		}
	}

	if node.isFull() {
		node1, node2, promoted = node.split()

		pointer, pointer1 = medium.Save(node1), medium.Save(node2)

		return
	}

	pointer = medium.Save(node)

	return
}

func (cursor *Cursor) Put(key, value []byte) (e error) {
	cursor.reset()

	cursor.offset, e = put(cursor.medium, cursor.offset, key, value)
	if e != nil {
		return
	}

	cursor.latest = key

	return
}
