package fifo

import (
	"github.com/voidDB/voidDB/common"
)

type FIFO []byte

func (fifo FIFO) headPointer() []byte {
	return common.Field(fifo, 0, wordSize)
}

func (fifo FIFO) getHeadPointer() int {
	return common.GetInt(
		fifo.headPointer(),
	)
}

func (fifo FIFO) setHeadPointer(pointer int) {
	common.PutInt(
		fifo.headPointer(),
		pointer,
	)

	return
}

func (fifo FIFO) nextIndex() []byte {
	return common.Field(fifo, wordSize, wordSize)
}

func (fifo FIFO) getNextIndex() int {
	return common.GetInt(
		fifo.nextIndex(),
	)
}

func (fifo FIFO) setNextIndex(pointer int) {
	common.PutInt(
		fifo.nextIndex(),
		pointer,
	)

	return
}

func (fifo FIFO) tailPointer() []byte {
	return common.Field(fifo, 2*wordSize, wordSize)
}

func (fifo FIFO) getTailPointer() int {
	return common.GetInt(
		fifo.tailPointer(),
	)
}

func (fifo FIFO) setTailPointer(pointer int) {
	common.PutInt(
		fifo.tailPointer(),
		pointer,
	)

	return
}
