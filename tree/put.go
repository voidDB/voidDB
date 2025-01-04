package tree

func Put(medium Medium, offset int, key, value []byte) (pointer int, e error) {
	var (
		newRoot  Node
		pointer1 int
		promoted []byte
	)

	pointer, e = medium.Save(value)
	if e != nil {
		return
	}

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

	pointer, e = medium.Save(newRoot)
	if e != nil {
		return
	}

	return
}

func _put(medium Medium, offset int, key []byte, putPointer, putValLen int) (
	pointer, pointer1 int, promoted []byte, e error,
) {
	var (
		node  Node = getNode(medium, offset, true)
		node1 Node
		node2 Node

		index  int
		valLen int
	)

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

		pointer, e = medium.Save(node1)
		if e != nil {
			return
		}

		pointer1, e = medium.Save(node2)
		if e != nil {
			return
		}

		return
	}

	pointer, e = medium.Save(node)
	if e != nil {
		return
	}

	return
}

func (cursor *Cursor) Put(key, value []byte) (e error) {
	cursor.reset()

	cursor.offset, e = Put(cursor.medium, cursor.offset, key, value)
	if e != nil {
		return
	}

	_, e = cursor._get(key)
	if e != nil {
		return
	}

	return
}
