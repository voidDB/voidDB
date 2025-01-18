package node

func copyElem(node1, node0 Node, index, shift int) {
	copy(
		node1.elem(index+shift),
		node0.elem(index),
	)

	return
}

func copyKey(node1, node0 Node, index, shift int) {
	copy(
		node1.key(index+shift),
		node0.key(index),
	)

	return
}

func copyElemKey(node1, node0 Node, index, shift int) {
	copyElem(node1, node0, index, shift)

	copyKey(node1, node0, index, shift)

	return
}
