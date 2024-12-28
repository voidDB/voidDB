package tree

func (node *Node) insert(index, pointer, valLen int, key []byte) (
	newNode Node,
) {
	var (
		i int
	)

	newNode = NewNode()

	for i = 0; i < index; i++ {
		copyNodeData(&newNode, node, i, 0)
	}

	newNode.setKey(i, key)

	newNode.setPointer(i, pointer)

	newNode.setValLen(i, valLen)

	for i++; i < node.length()+1; i++ {
		copyNodeData(&newNode, node, i, -1)
	}

	newNode.setPointer(i,
		node.pointer(i-1),
	)

	return
}
