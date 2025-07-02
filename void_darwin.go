//go:build darwin

// # Acknowledgements
//
// Andy Walker ([@flowchartsman]) quickly drew my attention to the fact that
// [syscall.Fdatasync] is undefined for GOOS=darwin, and saved me from
// accidentally breaking voidDB on macOS by replacing [os.File.Sync] with the
// unsupported syscall.
//
// [@flowchartsman]: https://github.com/flowchartsman
package voidDB

func (void *Void) fsync() error {
	return void.file.Sync()
}
