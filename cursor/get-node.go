package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func getNode(medium Medium, offset int, free bool) (
	n node.Node, dirty bool, e error,
) {
	n, dirty = medium.Load(offset, common.PageSize)

	e = n.VetMagic()
	if e != nil {
		return
	}

	if free {
		medium.Free(offset, common.PageSize)
	}

	return
}
