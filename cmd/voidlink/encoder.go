package main

import (
	"encoding/binary"
	"hash"
	"io"
	"sync"
	"syscall"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
)

const (
	maxMetadataLen = 7 * common.WordSize
	maxUintLen32   = 4

	offsetM = 0xd
	offsetX = 0xb
	offsetC = 0xa
)

// An Encoder is modelled after [encoding/gob.Encoder] from the Go standard
// library, but specialises in the transmission of key-value records.
//
// A record, consisting of a key no more than 512 bytes long, and a value of
// maximum size [math.MaxUint32] bytes, is encoded as follows:
//
//   - 2 bytes to represent the key length k in number of bytes,
//   - 1 <= x <= 4 bytes to represent the value length v in number of bytes,
//   - k bytes to hold the uninterpreted key,
//   - m * 8 bytes to hold uninterpreted, extra record metadata,
//   - v bytes to hold the uninterpreted value, and
//   - 4 bytes to hold an optional 32-bit checksum of the record.
//
// This incurs an overhead of 3 to 10 bytes per record, and leaves the first
// seven bits free to carry the following metadata:
//
//   - 3 bits to encode the value of m (from the fourth bullet point above),
//   - 2 bits to encode the value of x (from the third bullet point above), and
//   - 1 bit to indicate the presence of a trailing 32-bit checksum.
//
// Encoders are safe for concurrent use by multiple goroutines.
type Encoder struct {
	writer io.Writer
	hasher hash.Hash32

	mutex  sync.Mutex
	buffer []byte
}

// NewEncoder returns a new Encoder that will transmit on the [io.Writer], and
// optionally append a 32-bit checksum to every record if the [hash.Hash32] is
// not nil.
func NewEncoder(writer io.Writer, hasher hash.Hash32) (encoder *Encoder) {
	return &Encoder{
		writer: writer,
		hasher: hasher,
		buffer: make([]byte, maxUintLen32),
	}
}

// Encode transmits a key-value record.
func (encoder *Encoder) Encode(key, metadata, value []byte) (e error) {
	encoder.mutex.Lock()

	defer encoder.mutex.Unlock()

	switch {
	case len(key) > cursor.MaxKeyLength:
		fallthrough

	case len(metadata) > maxMetadataLen:
		fallthrough

	case len(metadata)%common.WordSize > 0:
		fallthrough

	case len(value) > cursor.MaxValueLength:
		return syscall.EINVAL
	}

	e = encoder.writeMXCK(key, metadata, value)
	if e != nil {
		return
	}

	e = encoder.writeV(value)
	if e != nil {
		return
	}

	e = encoder.writeAsIs(key)
	if e != nil {
		return
	}

	e = encoder.writeAsIs(metadata)
	if e != nil {
		return
	}

	e = encoder.writeAsIs(value)
	if e != nil {
		return
	}

	if encoder.hasher == nil {
		return
	}

	e = encoder.writeChecksum(key, metadata, value)
	if e != nil {
		return
	}

	return
}

func (encoder *Encoder) writeMXCK(key, metadata, value []byte) (e error) {
	// Writes the first two bytes, consisting of the following bit fields:
	//   * M: 4 bits to encode the record metadata length in number of 64-bit
	//        words,
	//   * X: 2 bits to encode the value of x, so that 1 <= x <= 4 represents
	//        len(value),
	//   * C: 1 bit to indicate the presence of a trailing 32-bit checksum,
	//   * K: 9 bits to represent len(key).
	//
	//  F E D C B A 9 8 7 6 5 4 3 2 1 0
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |  M  | X |C|         K         |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	var (
		m = uint16(len(metadata)/common.WordSize) << offsetM

		x = uint16(findX(value)%4) << offsetX
		// 1: 0b01, 2: 0b10, 3: 0b11, 4: 0b00

		c = uint16(1) << offsetC

		k = uint16(len(key))
	)

	if encoder.hasher == nil {
		c = 0
	}

	e = binary.Write(encoder.writer, binary.BigEndian, m|x|c|k)
	if e != nil {
		return
	}

	return
}

func (encoder *Encoder) writeV(value []byte) (e error) {
	// Writes one to four bytes representing len(value).

	binary.BigEndian.PutUint32(encoder.buffer,
		uint32(len(value)),
	)

	_, e = encoder.writer.Write(
		encoder.buffer[maxUintLen32-findX(value):],
	)
	if e != nil {
		return
	}

	return
}

func (encoder *Encoder) writeAsIs(bytes []byte) (e error) {
	// Writes the uninterpreted key or value.

	_, e = encoder.writer.Write(bytes)
	if e != nil {
		return
	}

	return
}

func (encoder *Encoder) writeChecksum(key, metadata, value []byte) (e error) {
	// Writes a 32-bit checksum of the record.

	defer encoder.hasher.Reset()

	_, e = encoder.hasher.Write(key)
	if e != nil {
		return
	}

	_, e = encoder.hasher.Write(metadata)
	if e != nil {
		return
	}

	_, e = encoder.hasher.Write(value)
	if e != nil {
		return
	}

	_, e = encoder.writer.Write(
		encoder.hasher.Sum(nil),
	)
	if e != nil {
		return
	}

	return
}

func findX(slice []byte) (x int) {
	// Returns the minimum number of bytes needed to encode an unsigned integer
	// indicating the length of slice.

	var (
		l int = len(slice)
	)

	switch {
	case l < 1<<8:
		return 1

	case l < 1<<16:
		return 2

	case l < 1<<24:
		return 3

	case l < 1<<32:
		return 4
	}

	return
}
