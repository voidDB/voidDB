package tree

func merge(node0, node1 Node) (newNode Node) {
	var (
		i int
	)

	newNode = NewNode()

	for i = 0; i < node0.length(); i++ {
		copyNodeData(&newNode, &node0, i, 0)
	}

	for i = i; i < MaxNodeLength-1; i++ {
		copyNodeData(&newNode, &node1, i,
			-node0.length(),
		)
	}

	newNode.setPointer(i,
		node1.pointer(
			i-node0.length(),
		),
	)

	return
}
