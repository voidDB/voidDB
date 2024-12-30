package tree

type Cursor struct {
	medium Medium

	offset int
	index  int

	stack []ancestor
}

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

func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	var (
		node Node = getNode(cursor.medium, cursor.offset, false)

		pointer int = node.pointer(cursor.index)
		valLen  int = node.valLen(cursor.index)
	)

	switch {
	case pointer == 0 && len(cursor.stack) == 0:
		e = ErrorNotFound

		return

	case pointer == 0:
		cursor.offset, cursor.index =
			cursor.stack[len(cursor.stack)-1].offset,
			cursor.stack[len(cursor.stack)-1].index

		cursor.stack = cursor.stack[:len(cursor.stack)-1]

		return cursor.GetNext()

	case valLen > 0:
		defer func() { cursor.index++ }()

		return node.key(cursor.index),
			cursor.medium.Load(pointer, valLen),
			nil
	}

	cursor.stack = append(cursor.stack,
		ancestor{cursor.offset, cursor.index + 1},
	)

	cursor.offset, cursor.index = pointer, 0

	return cursor.GetNext()
}

type ancestor struct {
	offset int
	index  int
}
