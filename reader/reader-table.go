package reader

import (
	"io"
	"math"
	"os"
	"syscall"

	"github.com/voidDB/voidDB/common"
)

const (
	maxNReaders = 1 << 22 // maximum number of PIDs allowed on most systems
	pathSuffix  = ".readers"
	wordSize    = common.WordSize
)

type ReaderTable struct {
	file *os.File
	mmap []byte

	locked map[int]struct{}
}

func NewReaderTable(path string) (e error) {
	var (
		file *os.File
	)

	file, e = os.Create(path + pathSuffix)
	if e != nil {
		return
	}

	e = file.Close()
	if e != nil {
		return
	}

	return
}

func OpenReaderTable(path string) (table *ReaderTable, e error) {
	table = &ReaderTable{
		locked: make(map[int]struct{}),
	}

	table.file, e = os.OpenFile(path+pathSuffix, os.O_RDWR, 0)
	if e != nil {
		return
	}

	table.mmap, e = syscall.Mmap(
		int(table.file.Fd()),
		0,
		maxNReaders*table.slotLength(),
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
	if e != nil {
		return
	}

	return
}

func (table *ReaderTable) AcquireSlot(txnID int) (
	releaseSlot func() error, e error,
) {
	var (
		index int
	)

	for index = 0; index < maxNReaders; index++ {
		releaseSlot, e = table.lockSlot(index)
		if e == nil {
			break
		}
	}

	e = table.setTxnID(index, txnID)
	if e != nil {
		return
	}

	return
}

func (table *ReaderTable) Close() (e error) {
	e = table.file.Close()
	if e != nil {
		return
	}

	e = syscall.Munmap(table.mmap)
	if e != nil {
		return
	}

	return
}

func (table *ReaderTable) OldestReader() (oldest int) {
	var (
		e     error
		index int
		size  int
		txnID int
	)

	size, e = table.fileSize()
	if e != nil {
		return // conservative; assumes oldest reader bears transaction ID 0
	}

	oldest = math.MaxInt64

	for index = 0; index < maxNReaders; index++ {
		if table.slotOffset(index) >= size {
			return
		}

		if table.slotIsLocked(index) { // reader is active
			txnID = table.getTxnID(index)

			if txnID < oldest {
				oldest = txnID
			}
		}
	}

	return
}

func (table *ReaderTable) fileSize() (size int, e error) {
	var (
		stat os.FileInfo
	)

	stat, e = table.file.Stat()
	if e != nil {
		return
	}

	size = int(stat.Size())

	return
}

func (table *ReaderTable) slotOffset(index int) int {
	return wordSize * index
}

func (table *ReaderTable) slotLength() int {
	return wordSize
}

func (table *ReaderTable) lockSlot(index int) (
	unlockSlot func() error, e error,
) {
	var (
		flock *syscall.Flock_t
	)

	if table.slotIsLocked(index) {
		return nil, syscall.EWOULDBLOCK
	}

	table.locked[index] = struct{}{}

	flock = &syscall.Flock_t{
		Type:   syscall.F_WRLCK,
		Whence: io.SeekStart,
		Start:  int64(table.slotOffset(index)),
		Len:    int64(table.slotLength()),
	}

	e = syscall.FcntlFlock(
		table.file.Fd(),
		fOFDSetlk, // process-indepedent; released when file desc. closed
		flock,
	)
	if e != nil {
		return
	}

	unlockSlot = func() error {
		defer delete(table.locked, index)

		flock.Type = syscall.F_UNLCK

		return syscall.FcntlFlock(
			table.file.Fd(),
			fOFDSetlk,
			flock,
		)
	}

	return
}

func (table *ReaderTable) slotIsLocked(index int) (locked bool) {
	var (
		flock *syscall.Flock_t
	)

	if _, locked = table.locked[index]; locked {
		return
	}

	flock = &syscall.Flock_t{
		Whence: io.SeekStart,
		Start:  int64(table.slotOffset(index)),
		Len:    int64(table.slotLength()),
	}

	syscall.FcntlFlock(
		table.file.Fd(),
		fOFDGetlk,
		flock,
	)

	return flock.Type != syscall.F_UNLCK
}

func (table *ReaderTable) getTxnID(index int) (txnID int) {
	return common.GetInt(
		common.Field(table.mmap,
			table.slotOffset(index),
			table.slotLength(),
		),
	)
}

func (table *ReaderTable) setTxnID(index, txnID int) (e error) {
	var (
		length int = table.slotLength()
		offset int = table.slotOffset(index)

		slot = make([]byte, length)
	)

	common.PutInt(slot, txnID)

	_, e = table.file.WriteAt(slot,
		int64(offset),
	)
	if e != nil {
		return
	}

	return
}
