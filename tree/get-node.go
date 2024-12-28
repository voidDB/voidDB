package tree

func getNode(medium Medium, offset int) Node {
	return Node(
		medium.Load(offset, PageSize),
	)
}
