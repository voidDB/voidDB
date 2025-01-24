package voidDB

import (
	"os"

	"golang.org/x/sys/unix"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
	"github.com/voidDB/voidDB/free"
	"github.com/voidDB/voidDB/reader"
)

type Txn struct {
	lockfile *os.File
	readers  *reader.ReaderTable

	read  readFunc
	write writeFunc
	sync  syncFunc

	meta     Meta
	saveList map[int][]byte
	freeList map[int][]int
	freeze   bool

	*cursor.Cursor
}

// Abort discards all changes made in a read-write transaction, and releases
// the exclusive write lock. In the case of a read-only transaction, Abort ends
// the moratorium on recycling of pages constituting its frozen-in-time view of
// the database. For this reason, applications should not be slow to abort
// transactions that have outlived their usefulness lest they prevent effective
// resource utilisation. Following an invocation of Abort, the transaction
// handle must no longer be used.
func (txn *Txn) Abort() (e error) {
	e = txn.readers.Close()
	if e != nil {
		return
	}

	if txn.lockfile == nil {
		goto end
	}

	e = txn.lockfile.Close()
	if e != nil {
		return
	}

end:
	*txn = Txn{}

	return
}

// Commit persists all changes to data made in a transaction. The state of the
// database is not really updated until Commit has been invoked. If it returns
// a nil error, effects of the transaction would be perceived in subsequent
// transactions, whereas pre-existing transactions will remain oblivious as
// intended. Whether Commit waits on [*os.File.Sync] depends on the mustSync
// argument passed to [*Void.BeginTxn]. The transaction handle is not safe to
// reuse after the first invocation of Commit, regardless of the result.
func (txn *Txn) Commit() (e error) {
	var (
		data   []byte
		offset int
	)

	txn.freeze = true

	txn.enqueueFreeList()

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

	if txn.sync != nil {
		e = txn.sync()
		if e != nil {
			return
		}
	}

	return txn.Abort()
}

func newTxn(path string, read readFunc, write writeFunc, sync syncFunc) (
	txn *Txn, e error,
) {
	txn = &Txn{
		read:     read,
		write:    write,
		sync:     sync,
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

	txn.readers, e = reader.OpenReaderTable(path)
	if e != nil {
		return
	}

	switch {
	case write == nil:
		txn.write = denyPermission

		e = txn.readers.AcquireSlot(
			txn.meta.getSerialNumber(),
		)
		if e != nil {
			return
		}

	default:
		txn.lockfile, e = os.OpenFile(path, os.O_RDONLY, 0)
		if e != nil {
			return
		}

		e = unix.Flock(
			int(txn.lockfile.Fd()),
			unix.LOCK_EX|unix.LOCK_NB,
		)
		if e != nil {
			return
		}
	}

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

func (txn *Txn) enqueueFreeList() {
	var (
		queue freeQueue

		head   int
		offset int
		size   int
		tail   int
	)

	for size = range txn.freeList {
		queue = txn.meta.freeQueue(size)

		head, tail = free.Put(medium{txn},
			queue.getTailPointer(),
			txn.meta.getSerialNumber(),
			txn.freeList[size],
		)

		if queue.getHeadPointer() == 0 {
			queue.setHeadPointer(head)
		}

		queue.setTailPointer(tail)

		for _, offset = range txn.freeList[size] {
			delete(txn.saveList, offset)
		}

		delete(txn.freeList, size)
	}

	return
}

type readFunc func(int, int) []byte

type writeFunc func([]byte, int) error

var (
	denyPermission writeFunc = func([]byte, int) error {
		return unix.EACCES
	}
)

type syncFunc func() error
