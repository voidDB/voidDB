package voidDB

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
	"github.com/voidDB/voidDB/reader"
)

// A Void is a handle on a database. To interact with the database, enter a
// transaction through [*Void.BeginTxn].
type Void struct {
	file *os.File
	mmap []byte
}

// NewVoid creates and initialises a database file and its reader table at path
// and path.readers respectively, and returns a handle on the database, or
// [os.ErrExist] if a file already exists at path. See also [OpenVoid] for an
// explanation of the capacity parameter.
func NewVoid(path string, capacity int) (void *Void, e error) {
	var (
		file *os.File
	)

	_, e = os.Stat(path)
	if e == nil {
		return nil, os.ErrExist
	}

	file, e = os.Create(path)
	if e != nil {
		return
	}

	defer file.Close()

	_, e = file.Write(
		newMetaInit(),
	)
	if e != nil {
		return
	}

	_, e = file.Write(
		newMetaInit(),
	)
	if e != nil {
		return
	}

	_, e = file.Write(
		node.NewNode(),
	)
	if e != nil {
		return
	}

	e = reader.NewReaderTable(path)
	if e != nil {
		return
	}

	return OpenVoid(path, capacity)
}

// OpenVoid returns a handle on the database persisted to the file at path.
//
// The capacity argument sets a hard limit on the size of the database file in
// number of bytes, but it applies only to transactions entered into via the
// database handle returned. The database file never shrinks, but it will not
// be allowed to grow if its size already exceeds capacity as at the time of
// invocation. A transaction running against the limit would incur
// [common.ErrorFull] on commit.
func OpenVoid(path string, capacity int) (void *Void, e error) {
	var (
		stat os.FileInfo
	)

	void = new(Void)

	void.file, e = os.OpenFile(path, os.O_RDWR, 0)
	if e != nil {
		return
	}

	stat, e = void.file.Stat()
	if e != nil {
		return
	}

	if int(stat.Size()) > capacity {
		capacity = int(stat.Size())
	}

	void.mmap, e = unix.Mmap(
		int(void.file.Fd()),
		0,
		capacity,
		unix.PROT_READ,
		unix.MAP_PRIVATE,
	)
	if e != nil {
		return
	}

	return
}

// View is similar to [*Void.Update], except that it begins and passes to
// operation a read-only transaction. If operation results in a non-nil error,
// that error is [errors.Join]-ed with the result of *Txn.Abort; otherwise only
// the latter is returned.
func (void *Void) View(operation func(*Txn) error) (e error) {
	var (
		txn *Txn
	)

	txn, e = void.BeginTxn(true, false)
	if e != nil {
		return
	}

	e = operation(txn)
	if e != nil {
		return errors.Join(e,
			txn.Abort(),
		)
	}

	return txn.Abort()
}

// Update is a convenient wrapper around [*Void.BeginTxn], [*Txn.Commit], and
// [*Txn.Abort], to help applications ensure timely termination of writers.
// If operation is successful (in that it returns a nil error), the transaction
// is automatically committed and the result of [*Txn.Commit] is returned.
// Otherwise, the transaction is aborted and the output of [errors.Join]
// wrapping the return values of operation and [*Txn.Abort] is returned.
// IMPORTANT: See also [*Void.BeginTxn] for an explanation of the mustSync
// parameter.
func (void *Void) Update(mustSync bool, operation func(*Txn) error) (e error) {
	var (
		txn *Txn
	)

	txn, e = void.BeginTxn(false, mustSync)
	if e != nil {
		return
	}

	e = operation(txn)
	if e != nil {
		return errors.Join(e,
			txn.Abort(),
		)
	}

	return txn.Commit()
}

// BeginTxn begins a new transaction. The resulting transaction cannot modify
// data if readonly is true: any changes made are isolated to the transaction
// and non-durable; otherwise it is a write transaction. Since there cannot be
// more than one ongoing write transaction per database at any point in time,
// the function may return [syscall.EAGAIN] or [syscall.EWOULDBLOCK] (same
// error, “resource temporarily unavailable”) if an uncommitted/unaborted
// incumbent is present in any thread/process in the system.
//
// Setting mustSync to true ensures that all changes to data are flushed to
// disk when the transaction is committed, at a cost to write performance;
// setting it to false empowers the filesystem to optimise writes at a risk of
// data loss in the event of a crash at the level of the operating system or
// lower, e.g. hardware or power failure. Database corruption is also
// conceivable, albeit only if the filesystem does not preserve write order.
// TL;DR: set mustSync to true if safety matters more than speed; false if vice
// versa.
//
// BeginTxn returns [common.ErrorResized] if the database file has grown beyond
// the capacity initially passed to [OpenVoid]. This can happen if another
// database handle with a higher capacity has been obtained via a separate
// invocation of OpenVoid in the meantime. To adapt to the new size and
// proceed, close the database handle and replace it with a new invocation of
// OpenVoid.
func (void *Void) BeginTxn(readonly, mustSync bool) (txn *Txn, e error) {
	var (
		stat  os.FileInfo
		sync  syncFunc
		write writeFunc
	)

	stat, e = void.file.Stat()
	if e != nil {
		return
	}

	if int(stat.Size()) > cap(void.mmap) {
		return nil, common.ErrorResized
	}

	if !readonly {
		write = void.write

		if mustSync {
			sync = void.file.Sync
		}
	}

	txn, e = newTxn(
		void.file.Name(),
		void.read,
		write,
		sync,
	)
	if e != nil {
		return
	}

	return
}

// Close closes the database file and releases the corresponding memory map,
// rendering both unusable to any remaining transactions already entered into
// using the database handle. These stranded transactions could give rise to
// undefined behaviour if their use is attempted, which could disrupt the
// application, but in any case they pose no danger whatsoever to the data
// safely jettisoned.
func (void *Void) Close() (e error) {
	e = void.file.Close()
	if e != nil {
		return
	}

	e = unix.Munmap(void.mmap)
	if e != nil {
		return
	}

	*void = Void{}

	return
}

func (void *Void) read(offset, length int) []byte {
	return void.mmap[offset : offset+length]
}

func (void *Void) write(data []byte, offset int) (e error) {
	if offset+len(data) > cap(void.mmap) {
		return common.ErrorFull
	}

	_, e = void.file.WriteAt(data,
		int64(offset),
	)
	if e != nil {
		return
	}

	return
}

func align(size int) int {
	return 1 << logarithm(size)
}

func logarithm(size int) (exp int) {
	for exp = 12; 1<<exp < size; exp++ {
		continue
	}

	return
}
