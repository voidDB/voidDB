//go:build darwin

// # Acknowledgements
//
// Had Nuno Cruces ([@ncruces], [u/ncruces]) not enlightened me to the fact that
// open file description (OFD) locks are supported by the XNU kernel, voidDB
// would most likely have remained unavailable on macOS.
//
// [@ncruces]: https://github.com/ncruces
// [u/ncruces]: https://www.reddit.com/user/ncruces/
package reader

const (
	// https://github.com/apple-oss-distributions/xnu/
	//   blob/8d741a5de7ff4191bf97d57b9f54c2f6d4a15585/
	//   bsd/sys/fcntl.h#L387-L389
	fOFDSetlk = 0x5a
	fOFDGetlk = 0x5c
)
