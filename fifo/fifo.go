package fifo

import (
	"github.com/voidDB/voidDB/common"
)

type FIFO []byte

func (fifo FIFO) headPointer() []byte {
	return common.WordN(fifo, 0)
}

func (fifo FIFO) getHeadPointer() int {
	return common.GetIntFromWord(
		fifo.headPointer(),
	)
}

func (fifo FIFO) setHeadPointer(pointer int) {
	common.PutIntIntoWord(
		fifo.headPointer(),
		pointer,
	)

	return
}

func (fifo FIFO) nextIndex() []byte {
	return common.WordN(fifo, 1)
}

func (fifo FIFO) getNextIndex() int {
	return common.GetIntFromWord(
		fifo.nextIndex(),
	)
}

func (fifo FIFO) setNextIndex(pointer int) {
	common.PutIntIntoWord(
		fifo.nextIndex(),
		pointer,
	)

	return
}

func (fifo FIFO) tailPointer() []byte {
	return common.WordN(fifo, 2)
}

func (fifo FIFO) getTailPointer() int {
	return common.GetIntFromWord(
		fifo.tailPointer(),
	)
}

func (fifo FIFO) setTailPointer(pointer int) {
	common.PutIntIntoWord(
		fifo.tailPointer(),
		pointer,
	)

	return
}
