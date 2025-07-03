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

		free Free
	)

	free, _ = medium.Page(offset)

	e = free.vetMagic()
	if e != nil {
		return
	}

	switch {
	case free.getTxnID() >= txnIDCeiling:
		fallthrough

	case free.getNextPointer() == 0:
		e = common.ErrorNotFound

		return

	case free.getLength() == 0:
		e = common.ErrorNotFound

		fallthrough

	case index+1 == free.getLength():
		medium.Free(offset,
			len(free),
		)

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
