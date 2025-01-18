package node

func (node Node) Update(index, pointer, length int) (newNode Node) {
	newNode = make([]byte, pageSize)

	copy(newNode, node)

	newNode.setValueOrChild(index, pointer, length)

	return
}
