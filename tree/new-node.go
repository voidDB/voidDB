package tree

func NewNode() (node Node) {
	node = make([]byte, pageSize)

	node.setMagic()

	return
}
