package free

import (
	"github.com/voidDB/voidDB/common"
)

func Get(medium Medium, offset, index int) (
	txnID, pointer, nextOffset, nextIndex int, e error,
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
		nextOffset, nextIndex = free.getNextPointer(), 0

	default:
		nextOffset, nextIndex = offset, index+1
	}

	txnID, pointer = free.getTxnID(), free.getPagePointer(index)

	return
}
