package node

func (node Node) Update(index, pointer, length int, metadata []byte) (
	newNode Node,
) {
	newNode = make([]byte, pageSize)

	copy(newNode, node)

	switch {
	case length < 0:
		newNode.elem(index).setPointer(pointer)

		newNode.elem(index).setExtraMetadata(metadata)

	default:
		newNode.setValueOrChild(index, pointer, length, metadata)
	}

	return
}
