package voidDB

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
	"github.com/voidDB/voidDB/node"
	"github.com/voidDB/voidDB/reader"
)

// A Txn is a transaction handle necessary for interacting with a database. A
// transaction is the sum of the state of the database as at the beginning of
// that transaction and any changes made within it. See [*Void.BeginTxn] for
// more information.
type Txn struct {
	lockfile *os.File
	readers  *reader.ReaderTable

	read  readFunc
	write writeFunc
	punch punchFunc
	sync  syncFunc

	meta     voidMeta
	saveList map[int][]byte
	freeWarm map[int][]int
	freeCool map[int][]int
	coolList map[int]struct{}
	freeze   bool

	*cursor.Cursor
}

// SerialNumber returns a serial number identifying a particular state of the
// database as at the beginning of the transaction. All transactions beginning
// from the same state share the same serial number.
func (txn *Txn) SerialNumber() int {
	return txn.meta.getSerialNumber()
}

// Timestamp returns the time as at the beginning of the transaction.
func (txn *Txn) Timestamp() time.Time {
	return txn.meta.getTimestamp()
}

// OpenCursor returns a handle on a cursor associated with the transaction and
// a particular keyspace. Keyspaces allow multiple datasets with potentially
// intersecting (overlapping) sets of keys to reside within the same database
// without conflict, provided that all keys are unique within their respective
// keyspaces. Argument keyspace must not be simultaneously non-nil and of zero
// length, or otherwise longer than [node.MaxKeyLength]. Passing nil as the
// argument causes OpenCursor to return a cursor in the default keyspace.
//
// CAUTION: An application utilising keyspaces should avoid modifying records
// within the default keyspace, as it is used to store pointers to all the
// other keyspaces. There is virtually no limit on the number of keyspaces in a
// database.
//
// Unless multiple keyspaces are required, there is usually no need to invoke
// OpenCursor because the transaction handle embeds a [*cursor.Cursor]
// associated with the default keyspace.
func (txn *Txn) OpenCursor(keyspace []byte) (c *cursor.Cursor, e error) {
	var (
		pointer []byte
	)

	if keyspace == nil {
		return txn.Cursor, nil
	}

	pointer, e = txn.Cursor.Get(keyspace)

	switch {
	case errors.Is(e, common.ErrorNotFound):
		e = txn.setRootNodePointer(keyspace,
			medium{txn, nil}.Save(
				node.NewNode(),
			),
		)
		if e != nil {
			return
		}

		return txn.OpenCursor(keyspace)

	case e != nil:
		return
	}

	c = cursor.NewCursor(medium{txn, keyspace},
		common.GetInt(pointer),
	)

	return
}

// Abort discards all changes made in a read-write transaction, and releases
// the exclusive write lock. In the case of a read-only transaction, Abort ends
// the moratorium on recycling of pages constituting its view of the dataset.
// For this reason, applications should not be slow to abort transactions that
// have outlived their usefulness lest they prevent effective resource
// utilisation. Following an invocation of Abort, the transaction handle must
// no longer be used.
func (txn *Txn) Abort() (e error) {
	switch {
	case txn.lockfile != nil:
		e = txn.lockfile.Close()

		fallthrough

	default:
		e = errors.Join(e,
			txn.readers.Close(),
		)
	}

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

		abort = func() {
			e = errors.Join(e,
				txn.Abort(),
			)
		}
	)

	defer abort()

	txn.freeze = true

	txn.enqueueFreeList()

	for offset, data = range txn.saveList {
		e = txn.write(data, offset)
		if e != nil {
			return
		}
	}

	if txn.sync != nil {
		e = txn.sync()
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

	return
}

func newTxn(path string, read readFunc, write writeFunc,
	punch punchFunc, seek seekFunc, sync syncFunc,
) (
	txn *Txn, e error,
) {
	txn = &Txn{
		read:     read,
		write:    write,
		punch:    punch,
		sync:     sync,
		saveList: make(map[int][]byte),
		freeWarm: make(map[int][]int),
		freeCool: make(map[int][]int),
		coolList: make(map[int]struct{}),
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

		e = syscall.Flock(
			int(txn.lockfile.Fd()),
			syscall.LOCK_EX|syscall.LOCK_NB,
		)
		if e != nil {
			return
		}
	}

	if seek != nil {
		txn.uncoverHoles(seek)
	}

	txn.Cursor = cursor.NewCursor(medium{txn, nil},
		txn.meta.getRootNodePointer(),
	)

	return
}

func (txn *Txn) setRootNodePointer(keyspace []byte, pointer int) (e error) {
	var (
		value []byte
	)

	if keyspace == nil {
		txn.meta.setRootNodePointer(pointer)

		return
	}

	value = make([]byte, wordSize)

	common.PutInt(value, pointer)

	return txn.Cursor.Put(keyspace, value)
}

func (txn *Txn) getMeta() (e error) {
	var (
		meta0 voidMeta = txn.read(0, pageSize)
		meta1 voidMeta = txn.read(pageSize, pageSize)
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
	txn.meta.setChecksum()

	return txn.write(txn.meta,
		txn.meta.getSerialNumber()%2*pageSize,
	)
}

func (txn *Txn) enqueueFreeList() {
	var (
		pointer int
		size    int
	)

	for size = range txn.freeWarm {
		txn.meta.freeQueue(size).Enqueue(
			medium{txn, nil},
			txn.meta.getSerialNumber(),
			txn.freeWarm[size],
		)
	}

	for size = range txn.freeCool {
		for _, pointer = range txn.freeCool[size] {
			txn.punch(pointer, size)
		}
	}

	return
}

func (txn *Txn) uncoverHoles(seek seekFunc) {
	var (
		e       error
		length  int
		pointer int
	)

	for {
		pointer, length, e = seek(pointer)
		if e != nil {
			length = txn.meta.getFrontierPointer() - pointer
		}

		for length > 0 {
			txn.freeCool[pageSize] = append(txn.freeCool[pageSize], pointer)

			txn.coolList[pointer] = struct{}{}

			pointer += pageSize

			length -= pageSize
		}

		if e != nil {
			break
		}
	}

	return
}

type readFunc func(int, int) []byte

type writeFunc func([]byte, int) error

var (
	denyPermission writeFunc = func([]byte, int) error {
		return syscall.EACCES
	}
)

type punchFunc func(int, int) error

type seekFunc func(int) (int, int, error)

type syncFunc func() error
