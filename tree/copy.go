package tree

func copyNodeData(node1, node0 *Node, index, shift int) {
	copy(
		node1._key(index),
		node0._key(index+shift),
	)

	copy(
		node1._keyLen(index),
		node0._keyLen(index+shift),
	)

	copy(
		node1._pointer(index),
		node0._pointer(index+shift),
	)

	copy(
		node1._valLen(index),
		node0._valLen(index+shift),
	)

	return
}
