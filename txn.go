package voidDB

import (
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
)

type Txn struct {
	read  readFunc
	write writeFunc
	quit  func() error

	meta     Meta
	saveList map[int][]byte
	freeList map[int]int

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
		offset int
	)

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
	)

	pointer = txn.meta.getFrontierPointer() // TODO: reuse free space

	txn.saveList[pointer] = make([]byte, length)

	copy(txn.saveList[pointer], data)

	txn.meta.setRootNodePointer(pointer)

	txn.meta.setFrontierPointer(pointer + length) // TODO: reuse free space

	return
}

func (txn medium) Free(offset, length int) {
	// delete(txn.saveList, offset) // FIXME

	// TODO: make free space available for reuse

	return
}

type readFunc func(int, int) []byte

type writeFunc func([]byte, int) error
