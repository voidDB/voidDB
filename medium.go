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

	pointer = txn.getFreePagePointer(length)

	txn.saveList[pointer] = data

	if !txn.freeze {
		txn.setRootNodePointer(txn.keyspace, pointer)
	}

	return
}

func (txn medium) SaveAt(offset int, data []byte) {
	txn.saveList[offset] = data

	return
}

func (txn medium) Free(offset, length int) {
	var (
		cool bool
	)

	length = align(length)

	switch _, cool = txn.coolList[offset]; cool {
	case true:
		txn.freeCool[length] = append(txn.freeCool[length], offset)

	default:
		txn.freeWarm[length] = append(txn.freeWarm[length], offset)
	}

	delete(txn.saveList, offset)

	return
}

func (txn medium) getFreePagePointer(size int) (pointer int) {
	var (
		e error
	)

	pointer = txn.getFreePageCool(size)

	if pointer > 0 {
		return
	}

	pointer, e = txn.getFreePageCold(size)

	if e != nil {
		pointer = txn.getFreePageNew(size)
	}

	txn.coolList[pointer] = struct{}{}

	return
}

func (txn medium) getFreePageCool(size int) (pointer int) {
	var (
		available bool
		pointers  []int
	)

	pointers, available = txn.freeCool[size]

	if !available {
		return -pageSize
	}

	pointer = pointers[0]

	txn.freeCool[size] = pointers[1:]

	if len(txn.freeCool[size]) == 0 {
		delete(txn.freeCool, size)
	}

	return
}

func (txn medium) getFreePageCold(size int) (pointer int, e error) {
	var (
		queue fifo.FIFO = txn.meta.freeQueue(size)
	)

	return queue.Dequeue(txn, txn.readers.OldestTxn)
}

func (txn medium) getFreePageNew(size int) (pointer int) {
	pointer = txn.meta.getFrontierPointer()

	txn.meta.setFrontierPointer(pointer + size)

	return
}
