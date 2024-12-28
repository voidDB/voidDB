package tree

func NewNode() (node Node) {
	node = make([]byte, PageSize)

	node.setMagic()

	return
}
