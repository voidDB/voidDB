package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func getNode(medium Medium, offset int, free bool) (
	n node.Node, e error,
) {
	n = medium.Load(offset, common.PageSize)

	if !n.IsNode() {
		e = common.ErrorCorrupt

		return
	}

	if free {
		medium.Free(offset, common.PageSize)
	}

	return
}
