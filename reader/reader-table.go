package reader

import (
	"errors"
	"io"
	"math"
	"os"
	"sync"
	"syscall"

	"github.com/voidDB/voidDB/common"
)

const (
	extension   = ".readers"
	maxNReaders = 1 << 22 // maximum number of PIDs allowed on most systems
	slotLength  = common.WordSize
)

type ReaderTable struct {
	file *os.File
	mmap []byte

	index int
	mutex sync.Mutex
	idPtr []*int
}

func NewReaderTable(path string) (e error) {

	var (
		file *os.File
	)

	file, e = os.Create(path + extension)
	if e != nil {
		return
	}

	return file.Close()
}

func OpenReaderTable(path string) (table *ReaderTable, e error) {
	table = new(ReaderTable)

	table.file, e = os.OpenFile(path+extension, os.O_RDWR, 0)
	if e != nil {
		return
	}

	table.mmap, e = syscall.Mmap(
		int(table.file.Fd()),
		0,
		maxNReaders*slotLength,
		syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if e != nil {
		return
	}

	table.index, e = table.acquireSlot()
	if e != nil {
		return
	}

	return table, table.setTxnIDPwrite(math.MaxInt64)
}

func (table *ReaderTable) acquireSlot() (index int, e error) {
	for index = 0; index < maxNReaders; index++ {
		e = table.lockSlot(index)
		if e == nil {
			return
		}
	}

	return
}

func (table *ReaderTable) Close() error {
	return errors.Join(
		table.file.Close(),
		syscall.Munmap(table.mmap),
	)
}

func (table *ReaderTable) AcquireHold(txnID int) (releaseHold func() error) {
	var (
		p *int = &txnID
	)

	table.mutex.Lock()

	defer table.mutex.Unlock()

	table.idPtr = append(table.idPtr, p)

	releaseHold = func() error {
		table.mutex.Lock()

		defer table.mutex.Unlock()

		*p = -1

		table.update()

		return nil
	}

	table.update()

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
		switch {
		case index*slotLength >= size:
			return

		case table.slotIsLocked(index): // reader is active
			txnID = table.getTxnID(index)

		case index == table.index && len(table.idPtr) > 0:
			txnID = *(table.idPtr[0])
		}

		if txnID < oldest {
			oldest = txnID
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

func (table *ReaderTable) lockSlot(index int) (e error) {
	var (
		flock *syscall.Flock_t
	)

	flock = &syscall.Flock_t{
		Type:   syscall.F_WRLCK,
		Whence: io.SeekStart,
		Start:  int64(slotLength * index),
		Len:    int64(slotLength),
	}

	e = syscall.FcntlFlock(
		table.file.Fd(),
		fOFDSetlk, // process-independent; released when file descriptor closed
		flock,
	)
	if e != nil {
		return
	}

	return
}

func (table *ReaderTable) slotIsLocked(index int) (locked bool) {
	var (
		flock = &syscall.Flock_t{
			Whence: io.SeekStart,
			Start:  int64(slotLength * index),
			Len:    int64(slotLength),
		}
	)

	syscall.FcntlFlock(
		table.file.Fd(),
		fOFDGetlk,
		flock,
	)

	return flock.Type != syscall.F_UNLCK
}

func (table *ReaderTable) update() {
	var (
		txnID int = math.MaxInt64
	)

	for {
		switch {
		case len(table.idPtr) == 0:

		case *table.idPtr[0] == -1:
			table.idPtr = table.idPtr[1:]

			continue

		default:
			txnID = *table.idPtr[0]
		}

		break
	}

	table.setTxnIDMemMap(txnID)

	return
}

func (table *ReaderTable) getTxnID(index int) (txnID int) {
	return common.GetIntFromWord(
		common.WordN(table.mmap, index),
	)
}

func (table *ReaderTable) setTxnIDMemMap(txnID int) {
	common.PutIntIntoWord(
		common.WordN(table.mmap, table.index),
		txnID,
	)

	return
}

func (table *ReaderTable) setTxnIDPwrite(txnID int) (e error) {
	var (
		word []byte = common.NewWord()
	)

	common.PutIntIntoWord(word, txnID)

	_, e = table.file.WriteAt(word,
		int64(table.index*slotLength),
	)
	if e != nil {
		return
	}

	return
}
