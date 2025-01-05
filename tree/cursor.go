package tree

type Cursor struct {
	medium Medium

	offset int
	index  int

	stack  []ancestor
	latest []byte
}

func (cursor *Cursor) reset() {
	if len(cursor.stack) > 0 {
		cursor.offset = cursor.stack[0].offset

		cursor.stack = cursor.stack[:0]
	}

	cursor.index = -1

	cursor.latest = nil

	return
}

func (cursor *Cursor) resume() {
	if cursor.latest == nil {
		return
	}

	cursor._get(cursor.latest)

	cursor.latest = nil

	return
}

type ancestor struct {
	offset int
	index  int
}
