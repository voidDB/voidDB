package cursor

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

func getNode(medium Medium, offset int, free bool) (
	n node.Node, e error,
) {
	if free {
		medium.Free(offset, common.PageSize)
	}

	n = node.Node(
		medium.Load(offset, common.PageSize),
	)

	if !n.IsNode() {
		e = common.ErrorCorrupt

		return
	}

	return
}
