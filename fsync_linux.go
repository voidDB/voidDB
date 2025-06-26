package voidDB

import (
	"syscall"
)

func (void *Void) fsync() error {
	return syscall.Fdatasync(
		int(void.file.Fd()),
	)
}
