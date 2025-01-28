package voidDB

import (
	"github.com/voidDB/voidDB/free"
)

type medium struct {
	*Txn

	keyspace []byte
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
	length = align(length)

	txn.freeList[length] = append(txn.freeList[length], offset)

	delete(txn.saveList, offset)

	return
}

func (txn medium) getFreePagePointer(size int) (pointer int) {
	var (
		available bool
		pointers  []int
	)

	if !txn.freeze {
		if pointers, available = txn.freeList[size]; available {
			for _, pointer = range pointers {
				txn.freeList[size] = txn.freeList[size][1:]

				if _, available = txn.freeSafe[pointer]; available {
					return
				}

				txn.freeList[size] = append(txn.freeList[size], pointer)
			}
		}
	}

	var (
		e          error
		nextIndex  int
		nextOffset int
		offset     int
		queue      freeQueue
		txnID      int
	)

	queue = txn.meta.freeQueue(size)

	offset = queue.getHeadPointer()

	txnID, pointer, nextOffset, nextIndex, e = free.Get(txn, offset,
		queue.getNextIndex(),
	)

	switch {
	case e != nil:
		fallthrough

	case txnID >= txn.readers.OldestTxn:
		pointer = txn.meta.getFrontierPointer()

		txn.meta.setFrontierPointer(pointer + size)

	case nextOffset != offset:
		txn.Free(offset, pageSize)

		queue.setHeadPointer(nextOffset)

		fallthrough

	default:
		queue.setNextIndex(nextIndex)
	}

	txn.freeSafe[pointer] = struct{}{}

	return
}
