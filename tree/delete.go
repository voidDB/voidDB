package tree

func (node *Node) delete(index int) (newNode Node) {
	var (
		i int
	)

	newNode = NewNode()

	for i = 0; i < index; i++ {
		copyNodeData(&newNode, node, i, 0)
	}

	for i = i; i < node.length(); i++ {
		copyNodeData(&newNode, node, i, 1)
	}

	newNode.setPointer(i,
		node.pointer(i+1),
	)

	return
}
