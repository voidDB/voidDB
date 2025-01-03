package tree

func (node *Node) split() (newNode0, newNode1 Node, promoted []byte) {
	var (
		i int
	)

	newNode0 = NewNode()
	newNode1 = NewNode()

	for i = 0; i < MaxNodeLength/2; i++ {
		copyNodeData(&newNode0, node, i, 0)
	}

	if node.valLen(i) == 0 {
		newNode0.setPointer(i,
			node.pointer(i),
		)

		promoted = node.key(i)

		defer func() {
			newNode1 = newNode1.delete(0)
		}()

	} else {
		promoted = node.key(i - 1)
	}

	for i = i; i < MaxNodeLength; i++ {
		copyNodeData(&newNode1, node, i-MaxNodeLength/2, MaxNodeLength/2)
	}

	newNode1.setPointer(i-MaxNodeLength/2,
		node.pointer(i),
	)

	return
}
