//go:build linux

package reader

const (
	// https://github.com/torvalds/linux/
	//   blob/7eb172143d5508b4da468ed59ee857c6e5e01da6/
	//   include/uapi/asm-generic/fcntl.h#L146-L147
	fOFDGetlk = 0x24
	fOFDSetlk = 0x25
)
