package voidDB

import (
	"os"
	"syscall"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
)

type Void struct {
	mmap []byte
	file *os.File
	size int
}

func NewVoid(path string, size int) (void *Void, e error) {
	var (
		file *os.File
	)

	size = pageAlign(size)

	file, e = os.Create(path)
	if e != nil {
		return
	}

	defer file.Close()

	e = syscall.Ftruncate(
		int(file.Fd()),
		int64(size),
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

	return OpenVoid(file.Name())
}

func OpenVoid(path string) (void *Void, e error) {
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

	void.size = int(stat.Size())

	void.mmap, e = syscall.Mmap(
		int(void.file.Fd()),
		0,
		void.size,
		syscall.PROT_READ,
		syscall.MAP_PRIVATE,
	)
	if e != nil {
		return
	}

	return
}

func (void *Void) BeginTxn(readonly bool) (txn *Txn, e error) {
	txn, e = newTxn(void.read, void.writeWhenReadonly)
	if e != nil {
		return
	}

	if readonly {
		return
	}

	e = void.lock()
	if e != nil {
		return
	}

	txn.write = void.write

	txn.quit = void.unlock

	return
}

func (void *Void) Close() (e error) {
	e = void.file.Close()
	if e != nil {
		return
	}

	e = syscall.Munmap(void.mmap)
	if e != nil {
		return
	}

	return
}

func Extend(path string, toSize int) (e error) {
	var (
		file *os.File
		stat os.FileInfo
	)

	toSize = pageAlign(toSize)

	file, e = os.Open(path)
	if e != nil {
		return
	}

	defer file.Close()

	stat, e = file.Stat()
	if e != nil {
		return
	}

	if stat.Size() > int64(toSize) {
		return
	}

	e = syscall.Ftruncate(
		int(file.Fd()),
		int64(toSize),
	)
	if e != nil {
		return
	}

	return
}

func (void *Void) lock() error {
	return syscall.Flock(
		int(void.file.Fd()),
		syscall.LOCK_EX|syscall.LOCK_NB,
	)
}

func (void *Void) unlock() error {
	return syscall.Flock(
		int(void.file.Fd()),
		syscall.LOCK_UN,
	)
}

func (void *Void) read(offset, length int) []byte {
	return void.mmap[offset : offset+length]
}

func (void *Void) write(data []byte, offset int) (e error) {
	var (
		length int = pageAlign(
			len(data),
		)
	)

	if offset+length > void.size {
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

func (void *Void) writeWhenReadonly([]byte, int) error {
	return syscall.EACCES
}

func pageAlign(size int) int {
	if size%pageSize > 0 {
		return size + pageSize - size%pageSize
	}

	return size
}
