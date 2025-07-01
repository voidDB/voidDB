package cursor

type position struct {
	offset int
	index  int
}

// A Cursor enables storage, retrieval, destruction of, and bidirectional
// iteration over key-value mappings in a keyspace via its methods. Out of the
// box, it is not safe for concurrent use; an application intending to do so
// should implement its own means of ensuring mutual exclusion.
type Cursor struct {
	medium Medium
	stack  []position
	latest []byte

	position
}

// NewCursor is a low-level constructor used by
// [*github.com/voidDB/voidDB.Txn.OpenCursor] and
// [*github.com/voidDB/voidDB.Void.BeginTxn].
func NewCursor(medium Medium, offset int) *Cursor {
	const (
		maxStackDepth = 26 // = log4(2^64 / 4096)
		// Assuming MaxNodeLength = 7 (such that each node must have at least
		// four children) and PageSize = 4096, a stack depth of 26 would put
		// within reach every node that could ever exist in a 64-bit address
		// space. 18 would suffice for the more common 48-bit address space.
	)

	return &Cursor{
		medium:   medium,
		stack:    make([]position, 0, maxStackDepth),
		position: position{offset, -1},
	}
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

	cursor.get(cursor.latest)

	cursor.latest = nil

	return
}

func (cursor *Cursor) descend(offset, index int) {
	cursor.stack = append(cursor.stack,
		position{cursor.offset, cursor.index},
	)

	cursor.position = position{offset, index}

	return
}

func (cursor *Cursor) ascend() {
	cursor.position = cursor.stack[len(cursor.stack)-1]

	cursor.stack = cursor.stack[:len(cursor.stack)-1]

	return
}
