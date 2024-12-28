package tree

func Del(medium Medium, offset int, key []byte) (pointer int, e error) {
	var (
		newRoot Node
	)

	newRoot, e = _del(medium, offset, key)
	if e != nil {
		return
	}

	pointer, e = medium.Save(newRoot)
	if e != nil {
		return
	}

	return
}

func _del(medium Medium, offset int, key []byte) (newNode Node, e error) {
	var (
		node Node = getNode(medium, offset)

		child  Node
		child1 Node

		index   int
		pointer int
		valLen  int
	)

	index, pointer, valLen = node.search(key)

	switch {
	case pointer == 0:
		e = ErrorNotFound

		return

	case valLen > 0:
		newNode = node.delete(index)

		return

	default:
		child, e = _del(medium, pointer, key)
		if e != nil {
			return
		}
	}

	if !child.isPoor() {
		pointer, e = medium.Save(child)
		if e != nil {
			return
		}

		newNode = node.update(index, pointer, 0)

		return
	}

	switch {
	case index > 0:
		child1 = getNode(medium,
			node.pointer(index-1),
		)

	default: // child is leftmost
		child1 = getNode(medium,
			node.pointer(1),
		)
	}

	switch {
	case index > 0 && child1.isRich():
		child = child.insert(0, 0, 0, nil)

		copyNodeData(&child, &child1, 0,
			child1.length()-1,
		)

		pointer, e = medium.Save(child)
		if e != nil {
			return
		}

		newNode = node.update(index, pointer, 0)

		child1 = child1.delete(
			child1.length() - 1,
		)

		newNode.setKey(index-1,
			child1.key(
				child1.length()-1,
			),
		)

		pointer, e = medium.Save(child1)
		if e != nil {
			return
		}

		newNode.setPointer(index-1, pointer)

	case index > 0:
		child = merge(child1, child)

		pointer, e = medium.Save(child)
		if e != nil {
			return
		}

		newNode = node.delete(index - 1)

		newNode.setPointer(index-1, pointer)

	case child1.isRich():
		copyNodeData(&child, &child1,
			child.length(),
			-child.length(),
		)

		pointer, e = medium.Save(child)
		if e != nil {
			return
		}

		newNode = node.update(0, pointer, 0)

		newNode.setKey(0,
			child1.key(0),
		)

		child1 = child1.delete(0)

		pointer, e = medium.Save(child1)
		if e != nil {
			return
		}

		newNode.setPointer(1, pointer)

	default:
		child = merge(child, child1)

		pointer, e = medium.Save(child)
		if e != nil {
			return
		}

		newNode = node.delete(0)

		newNode.setPointer(0, pointer)
	}

	if newNode.length() == 0 {
		newNode = child
	}

	return
}
