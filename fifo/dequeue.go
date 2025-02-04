package fifo

import (
	"github.com/voidDB/voidDB/common"
)

func (fifo FIFO) Dequeue(medium Medium, txnIDCeiling int) (
	pointer int, e error,
) {
	var (
		index  int = fifo.getNextIndex()
		offset int = fifo.getHeadPointer()

		free Free = medium.Load(offset, pageSize)
	)

	switch {
	case !free.isFree():
		e = common.ErrorCorrupt

		return

	case free.getTxnID() >= txnIDCeiling:
		fallthrough

	case free.getLength() == 0:
		e = common.ErrorNotFound

		return

	case index+1 == free.getLength():
		medium.Free(offset, pageSize)

		fifo.setHeadPointer(
			free.getNextPointer(),
		)

		fifo.setNextIndex(0)

	default:
		fifo.setNextIndex(index + 1)
	}

	pointer = free.getPagePointer(index)

	return
}
