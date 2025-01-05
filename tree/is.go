package tree

func (node *Node) isFull() bool {
	return node.length() == maxNodeLength
}

func (node *Node) isRich() bool {
	return node.length() > maxNodeLength/2
}

func (node *Node) isPoor() bool {
	return node.length() < maxNodeLength/2
}
