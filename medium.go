package voidDB

import (
	"github.com/voidDB/voidDB/fifo"
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
		e     error
		queue fifo.FIFO
	)

	queue = txn.meta.freeQueue(size)

	pointer, e = queue.Dequeue(txn, txn.readers.OldestTxn)
	if e != nil {
		pointer = txn.meta.getFrontierPointer()

		txn.meta.setFrontierPointer(pointer + size)
	}

	txn.freeSafe[pointer] = struct{}{}

	return
}
