package voidDB

import (
	"os"

	"golang.org/x/sys/unix"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/node"
	"github.com/voidDB/voidDB/reader"
)

type Void struct {
	file *os.File
	mmap []byte
}

func NewVoid(path string, capacity int) (void *Void, e error) {
	var (
		file *os.File
	)

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

func (void *Void) BeginTxn(readonly bool) (txn *Txn, e error) {
	var (
		write writeFunc
	)

	if !readonly {
		write = void.write
	}

	txn, e = newTxn(
		void.file.Name(),
		void.read,
		write,
	)
	if e != nil {
		return
	}

	return
}

func (void *Void) Close() (e error) {
	e = void.file.Close()
	if e != nil {
		return
	}

	e = unix.Munmap(void.mmap)
	if e != nil {
		return
	}

	return
}

func (void *Void) read(offset, length int) []byte {
	return void.mmap[offset : offset+length]
}

func (void *Void) write(data []byte, offset int) (e error) {
	_, e = void.file.WriteAt(data,
		int64(offset),
	)
	if e != nil {
		return
	}

	if offset+len(data) > cap(void.mmap) {
		return common.ErrorFull
	}

	return
}

func align(size int) int {
	return 1 << log(size)
}

func log(size int) (exp int) {
	for exp = 12; 1<<exp < size; exp++ {
		continue
	}

	return
}
