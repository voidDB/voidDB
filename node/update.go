package node

func (node Node) Update(index, pointer, length int, elemMeta []byte,
	inPlace bool,
) (
	newNode Node,
) {
	switch {
	case inPlace:
		newNode = node

	default:
		newNode = node.copy()
	}

	switch {
	case length < 0:
		newNode.elem(index).setPointer(pointer)

		newNode.elem(index).setMeta(elemMeta)

	default:
		newNode.setValueOrChild(index, pointer, length, elemMeta)
	}

	return
}
