package tree

func getNode(medium Medium, offset int, free bool) Node {
	if free {
		medium.Free(offset, pageSize)
	}

	return Node(
		medium.Load(offset, pageSize),
	)
}
