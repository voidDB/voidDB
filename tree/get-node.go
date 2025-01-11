package tree

func getNode(medium Medium, offset int, free bool) (node Node, e error) {
	if free {
		medium.Free(offset, pageSize)
	}

	node = Node(
		medium.Load(offset, pageSize),
	)

	if !node.isNode() {
		e = errorCorrupt

		return
	}

	return
}
