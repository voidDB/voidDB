package tree

type Cursor struct {
	medium Medium

	offset int
	index  int

	ancestors []int
	anIndices []int
}

func NewCursor(medium Medium, offset int) *Cursor {
	const (
		maxStackDepth = 512
	)

	return &Cursor{
		medium:    medium,
		offset:    offset,
		ancestors: make([]int, 0, maxStackDepth),
		anIndices: make([]int, 0, maxStackDepth),
	}
}

func (cursor *Cursor) GetNext() (key, value []byte, e error) {
	var (
		node Node = getNode(cursor.medium, cursor.offset, false)

		pointer int = node.pointer(cursor.index)
		valLen  int = node.valLen(cursor.index)
	)

	switch {
	case pointer == 0 && len(cursor.ancestors) == 0:
		e = ErrorNotFound

		return

	case pointer == 0:
		cursor.offset = cursor.ancestors[len(cursor.ancestors)-1]

		cursor.index = cursor.anIndices[len(cursor.anIndices)-1]

		cursor.ancestors = cursor.ancestors[:len(cursor.ancestors)-1]

		cursor.anIndices = cursor.anIndices[:len(cursor.anIndices)-1]

		return cursor.GetNext()

	case valLen > 0:
		defer func() { cursor.index++ }()

		return node.key(cursor.index),
			cursor.medium.Load(pointer, valLen),
			nil

	default:
		cursor.ancestors = append(cursor.ancestors, cursor.offset)

		cursor.anIndices = append(cursor.anIndices, cursor.index+1)

		cursor.offset = pointer

		cursor.index = 0

		return cursor.GetNext()
	}

	return
}
