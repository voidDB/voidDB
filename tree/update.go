package tree

func (node *Node) update(index, pointer, valLen int) (newNode Node) {
	var (
		i int
	)

	newNode = NewNode()

	for i = 0; i <= index; i++ {
		copyNodeData(&newNode, node, i, 0)
	}

	newNode.setPointer(index, pointer)

	newNode.setValLen(index, valLen)

	for i = i; i < node.length(); i++ {
		copyNodeData(&newNode, node, i, 0)
	}

	newNode.setPointer(i,
		node.pointer(i),
	)

	return
}
