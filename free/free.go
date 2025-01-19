package free

import (
	"bytes"

	"github.com/voidDB/voidDB/common"
)

const (
	MaxNodeLength = 508

	pageSize = common.PageSize
	wordSize = common.WordSize

	freeMagic = "voidFREE"
)

type Free []byte

func NewFree(txnID int) (free Free) {
	free = make([]byte, pageSize)

	free.setMagic()

	free.setTxnID(txnID)

	return
}

func (free Free) magic() []byte {
	return common.Field(free, 0, wordSize)
}

func (free Free) isFree() bool {
	return bytes.Equal(
		free.magic(),
		[]byte(freeMagic),
	)
}

func (free Free) setMagic() {
	copy(
		free.magic(),
		[]byte(freeMagic),
	)

	return
}

func (free Free) txnID() []byte {
	return common.Field(free, wordSize, wordSize)
}

func (free Free) getTxnID() int {
	return common.GetInt(
		free.txnID(),
	)
}

func (free Free) setTxnID(txnID int) {
	common.PutInt(
		free.txnID(),
		txnID,
	)

	return
}

func (free Free) length() []byte {
	return common.Field(free, 2*wordSize, wordSize)
}

func (free Free) getLength() int {
	return common.GetInt(
		free.length(),
	)
}

func (free Free) setLength(length int) {
	common.PutInt(
		free.length(),
		length,
	)

	return
}

func (free Free) nextPointer() []byte {
	return common.Field(free, 3*wordSize, wordSize)
}

func (free Free) getNextPointer() int {
	return common.GetInt(
		free.nextPointer(),
	)
}

func (free Free) setNextPointer(pointer int) {
	common.PutInt(
		free.nextPointer(),
		pointer,
	)

	return
}

func (free Free) pagePointer(index int) []byte {
	return common.Field(free,
		(4+index)*wordSize,
		wordSize,
	)
}

func (free Free) getPagePointer(index int) int {
	return common.GetInt(
		free.pagePointer(index),
	)
}

func (free Free) setPagePointer(index, pointer int) {
	common.PutInt(
		free.pagePointer(index),
		pointer,
	)

	return
}
