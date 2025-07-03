package fifo

import (
	"github.com/voidDB/voidDB/common"
)

const (
	MaxNodeLength = 508
)

var (
	freeMagic = []byte("voidFREE")
)

type Free []byte

func NewFree(txnID int) (free Free) {
	free = common.NewPage()

	free.setMagic()

	free.setTxnID(txnID)

	return
}

func (free Free) magic() []byte {
	return common.WordN(free, 0)
}

func (free Free) setMagic() {
	copy(
		free.magic(),
		freeMagic,
	)

	return
}

func (free Free) vetMagic() error {
	return common.ErrorIfNotEqual(
		free.magic(),
		freeMagic,
		common.ErrorCorrupt,
	)
}

func (free Free) txnID() []byte {
	return common.WordN(free, 1)
}

func (free Free) getTxnID() int {
	return common.GetIntFromWord(
		free.txnID(),
	)
}

func (free Free) setTxnID(txnID int) {
	common.PutIntIntoWord(
		free.txnID(),
		txnID,
	)

	return
}

func (free Free) length() []byte {
	return common.WordN(free, 2)
}

func (free Free) getLength() int {
	return common.GetIntFromWord(
		free.length(),
	)
}

func (free Free) setLength(length int) {
	common.PutIntIntoWord(
		free.length(),
		length,
	)

	return
}

func (free Free) nextPointer() []byte {
	return common.WordN(free, 3)
}

func (free Free) getNextPointer() int {
	return common.GetIntFromWord(
		free.nextPointer(),
	)
}

func (free Free) setNextPointer(pointer int) {
	common.PutIntIntoWord(
		free.nextPointer(),
		pointer,
	)

	return
}

func (free Free) pagePointer(index int) []byte {
	return common.WordN(free, 4+index)
}

func (free Free) getPagePointer(index int) int {
	return common.GetIntFromWord(
		free.pagePointer(index),
	)
}

func (free Free) setPagePointer(index, pointer int) {
	common.PutIntIntoWord(
		free.pagePointer(index),
		pointer,
	)

	return
}
