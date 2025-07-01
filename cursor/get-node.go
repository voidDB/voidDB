package cursor

import (
	"github.com/voidDB/voidDB/node"
)

func getNode(medium Medium, offset int, free bool) (
	n node.Node, dirty bool, e error,
) {
	n, dirty = medium.Page(offset)

	e = n.VetMagic()
	if e != nil {
		return
	}

	if free {
		medium.Free(offset,
			len(n),
		)
	}

	return
}
