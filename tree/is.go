package tree

func (node *Node) isFull() bool {
	return node.length() == MaxNodeLength
}

func (node *Node) isRich() bool {
	return node.length() > MaxNodeLength/2
}

func (node *Node) isPoor() bool {
	return node.length() < MaxNodeLength/2
}
