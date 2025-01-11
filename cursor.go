package voidDB

import (
	"github.com/voidDB/voidDB/tree"
)

type Cursor struct {
	*tree.Cursor
}

func newCursor(txn *Txn) Cursor {
	return Cursor{
		tree.NewCursor(medium{txn},
			txn.meta.rootNodePointer(),
		),
	}
}
