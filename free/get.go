package free

import (
	"github.com/voidDB/voidDB/common"
)

func Get(medium Medium, offset, index int) (
	pointer, nextOffset, nextIndex int, e error,
) {
	var (
		free Free = medium.Load(offset, pageSize)
	)

	switch {
	case !free.isFree():
		e = common.ErrorCorrupt

		return

	case free.getLength() == 0:
		e = common.ErrorNotFound

		return

	case index+1 == free.getLength():
		medium.Free(offset, pageSize)

		nextOffset = free.getNextPointer()

		nextIndex = 0

	default:
		nextOffset = offset

		nextIndex = index + 1
	}

	pointer = free.getPagePointer(index)

	return
}
