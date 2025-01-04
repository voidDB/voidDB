package tree

type Cursor struct {
	medium Medium

	offset int
	index  int

	stack []ancestor
}

func (cursor *Cursor) reset() {
	if len(cursor.stack) > 0 {
		cursor.offset = cursor.stack[0].offset

		cursor.stack = cursor.stack[:0]
	}

	cursor.index = -1

	return
}

type ancestor struct {
	offset int
	index  int
}
