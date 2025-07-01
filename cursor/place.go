package cursor

import (
	"github.com/voidDB/voidDB/node"
)

type place int

const (
	First place = 1
	Last        = -1
)

func (cursor *Cursor) Place(where place) {
	cursor.reset()

	switch where {
	case First:

	case Last:
		cursor.index = node.MaxNodeLength
	}

	return
}
