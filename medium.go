package voidDB

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/fifo"
)

type medium struct {
	*Txn

	keyspace []byte
}

func (txn medium) Meta() []byte {
	return common.Field(txn.meta, 2*wordSize, 2*wordSize)
}

func (txn medium) Load(offset, length int) (data []byte) {
	var (
		cached bool
	)

	data, cached = txn.saveList[offset]

	if cached {
		return data[:length]
	}

	return txn.read(offset, length)
}

func (txn medium) Save(data []byte) (pointer int) {
	var (
		length int = align(
			len(data),
		)
	)

	switch txn.freeze {
	case true:
		pointer = txn.getNewPagePointer(length)

	default:
		pointer = txn.getFreePagePointer(length)

		txn.setRootNodePointer(txn.keyspace, pointer)
	}

	txn.saveList[pointer] = data

	return
}

func (txn medium) SaveAt(offset int, data []byte) {
	txn.saveList[offset] = data

	return
}

func (txn medium) Free(offset, length int) {
	var (
		warm bool
	)

	length = align(length)

	switch _, warm = txn.warmList[offset]; warm {
	case true:
		txn.freeWarm[length] = append(txn.freeWarm[length], offset)

	default:
		txn.freeCold[length] = append(txn.freeCold[length], offset)
	}

	delete(txn.saveList, offset)

	return
}

func (txn medium) getFreePagePointer(size int) (pointer int) {
	var (
		e error
	)

	pointer = txn.getFreePageWarm(size)

	if pointer > 0 {
		return
	}

	pointer, e = txn.getFreePageCold(size)

	if e != nil {
		pointer = txn.getNewPagePointer(size)
	}

	txn.warmList[pointer] = struct{}{}

	return
}

func (txn medium) getFreePageWarm(size int) (pointer int) {
	var (
		available bool
		pointers  []int
	)

	pointers, available = txn.freeWarm[size]

	if !available || len(pointers) == 0 {
		return
	}

	pointer = pointers[0]

	txn.freeWarm[size] = pointers[1:]

	return
}

func (txn medium) getFreePageCold(size int) (pointer int, e error) {
	var (
		queue fifo.FIFO = txn.meta.freeQueue(size)
	)

	return queue.Dequeue(txn, txn.readers.OldestTxn)
}

func (txn medium) getNewPagePointer(size int) (pointer int) {
	pointer = txn.meta.getFrontierPointer()

	txn.meta.setFrontierPointer(pointer + size)

	return
}
