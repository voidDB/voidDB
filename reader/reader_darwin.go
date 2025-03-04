//go:build darwin

package reader

const (
	// https://github.com/apple-oss-distributions/xnu/
	//   blob/8d741a5de7ff4191bf97d57b9f54c2f6d4a15585/
	//   bsd/sys/fcntl.h#L387-L389
	fOFDSetlk = 0x5a
	fOFDGetlk = 0x5c
)
