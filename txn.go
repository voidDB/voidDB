package voidDB

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
	"github.com/voidDB/voidDB/free"
)

type Txn struct {
	read  readFunc
	write writeFunc
	quit  func() error

	meta     Meta
	saveList map[int][]byte
	freeList map[int][]int
	freeze   bool

	*cursor.Cursor
}

func (txn *Txn) Abort() (e error) {
	e = txn.quit()
	if e != nil {
		return
	}

	txn = nil

	return
}

func (txn *Txn) Commit() (e error) {
	var (
		data   []byte
		head   int
		offset int
		size   int
		tail   int
	)

	txn.freeze = true

	for size = range txn.freeList {
		head, tail = free.Put(medium{txn},
			txn.meta.getFreeListTailPtr(size),
			txn.meta.getSerialNumber(),
			txn.freeList[size],
		)

		if txn.meta.getFreeListHeadPtr(size) == 0 {
			txn.meta.setFreeListHeadPtr(size, head)
		}

		txn.meta.setFreeListTailPtr(size, tail)

		for _, offset = range txn.freeList[size] {
			delete(txn.saveList, offset)
		}
	}

	for offset, data = range txn.saveList {
		e = txn.write(data, offset)
		if e != nil {
			return
		}
	}

	e = txn.putMeta()
	if e != nil {
		return
	}

	return txn.Abort()
}

func newTxn(read readFunc, write writeFunc) (txn *Txn, e error) {
	txn = &Txn{
		read:     read,
		write:    write,
		saveList: make(map[int][]byte),
		freeList: make(map[int][]int),
	}

	e = txn.getMeta()
	if e != nil {
		return
	}

	txn.meta.setTimestamp()

	txn.meta.setSerialNumber(
		txn.meta.getSerialNumber() + 1,
	)

	txn.Cursor = cursor.NewCursor(medium{txn},
		txn.meta.getRootNodePointer(),
	)

	return
}

func (txn *Txn) getMeta() (e error) {
	var (
		meta0 Meta = Meta(txn.read(0, pageSize))
		meta1 Meta = Meta(txn.read(pageSize, pageSize))
	)

	switch {
	case meta0.isMeta() && meta1.isMeta() &&
		meta0.getSerialNumber() < meta1.getSerialNumber():
		txn.meta = meta1.makeCopy()

	case meta0.isMeta() && meta1.isMeta():
		txn.meta = meta0.makeCopy()

	case meta0.isMeta():
		txn.meta = meta0.makeCopy()

	case meta1.isMeta():
		txn.meta = meta1.makeCopy()

	default:
		e = common.ErrorInvalid
	}

	return
}

func (txn *Txn) putMeta() error {
	return txn.write(txn.meta,
		txn.meta.getSerialNumber()%2*pageSize,
	)
}

type medium struct {
	*Txn
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
		length int = pageAlign(
			len(data),
		)

		e error
	)

	pointer, e = txn.getFreePagePointer(length)
	if e != nil {
		pointer = txn.meta.getFrontierPointer()

		txn.meta.setFrontierPointer(pointer + length)
	}

	txn.saveList[pointer] = make([]byte, length)

	copy(txn.saveList[pointer], data)

	if !txn.freeze {
		txn.meta.setRootNodePointer(pointer)
	}

	return
}

func (txn medium) SaveAt(offset int, data []byte) {
	txn.saveList[offset] = data

	return
}

func (txn medium) Free(offset, length int) {
	length = pageAlign(length) // FIXME

	txn.freeList[length] = append(txn.freeList[length], offset)

	return
}

func (txn medium) getFreePagePointer(size int) (pointer int, e error) {
	var (
		nextOffset int
		nextIndex  int
	)

	pointer, nextOffset, nextIndex, e = free.Get(txn,
		txn.meta.getFreeListHeadPtr(size),
		txn.meta.getFreeListNextIdx(size),
	)
	if e != nil {
		return
	}

	// TODO: check reader table!

	txn.meta.setFreeListHeadPtr(size, nextOffset)

	txn.meta.setFreeListNextIdx(size, nextIndex)

	return
}

type readFunc func(int, int) []byte

type writeFunc func([]byte, int) error
