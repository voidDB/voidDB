//go:build !linux

package voidDB

func (void *Void) fsync() error {
	return void.file.Sync()
}
