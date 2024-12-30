package tree

func NewCursor(medium Medium, offset int) *Cursor {
	const (
		maxStackDepth = 26 // = log4(2^64 / 4096)
		// Assuming MaxNodeLength = 7 (such that each node must have at least
		// four children) and PageSize = 4096, a stack depth of 26 would put
		// within reach every node that could ever exist in a 64-bit address
		// space. 18 would suffice for the more common 48-bit address space.
	)

	return &Cursor{
		medium: medium,
		offset: offset,
		stack:  make([]ancestor, 0, maxStackDepth),
	}
}
