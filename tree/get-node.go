package tree

func getNode(medium Medium, offset int, free bool) Node {
	if free {
		medium.Free(offset, PageSize)
	}

	return Node(
		medium.Load(offset, PageSize),
	)
}
