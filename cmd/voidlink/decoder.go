package main

import (
	"encoding/binary"
	"hash"
	"io"
	"sync"
	"syscall"

	"github.com/voidDB/voidDB/common"
)

// Inspired by [encoding/gob.Decoder] from the Go standard library, a Decoder
// specialises in the receipt of key-value records transmitted by an Encoder
// counterpart. It is safe for concurrent use by multiple goroutines.
type Decoder struct {
	reader io.Reader
	hasher hash.Hash32

	mutex  sync.Mutex
	buffer []byte
}

// NewDecoder returns a new Decoder that will receive from the [io.Reader], and
// optionally verify the checksum of every record if the [hash.Hash32] is not
// nil.
func NewDecoder(reader io.Reader, hasher hash.Hash32) (decoder *Decoder) {
	return &Decoder{
		reader: reader,
		hasher: hasher,
		buffer: make([]byte, maxUintLen32),
	}
}

// Decode receives the next record from the input stream and returns three byte
// slices containing the record timestamp, key and value, respectively.
//
// At the end of the stream, Decode returns an [io.EOF].
func (decoder *Decoder) Decode() (key, metadata, value []byte, e error) {
	var (
		c bool // a trailing 32-bit checksum is present if true
		k int  // key length
		m int  // metadata length
		v int  // value length
		x int  // number of bytes representing value length
	)

	decoder.mutex.Lock()

	defer decoder.mutex.Unlock()

	m, x, c, k, e = decoder.readMXCK()
	if e != nil {
		return
	}

	v, e = decoder.readV(x)
	if e != nil {
		return
	}

	key, e = decoder.readAsIs(k)
	if e != nil {
		return
	}

	metadata, e = decoder.readAsIs(m)
	if e != nil {
		return
	}

	value, e = decoder.readAsIs(v)
	if e != nil {
		return
	}

	if !c {
		return
	}

	e = decoder.verifyChecksum(key, metadata, value)
	if e != nil {
		return
	}

	return
}

func (decoder *Decoder) readAsIs(n int) (bytes []byte, e error) {
	// Reads n bytes uninterpreted.

	bytes = make([]byte, n)

	_, e = io.ReadFull(decoder.reader, bytes)
	if e != nil {
		return
	}

	return
}

func (decoder *Decoder) readMXCK() (m int, x int, c bool, k int, e error) {
	// Reads the first two bytes, expecting the following bit fields:
	//   * M: 3 bits to encode the record metadata length in number of 64-bit
	//        words,
	//   * X: 2 bits to encode the value of x, so that 1 <= x <= 4 represents
	//        len(value),
	//   * C: 1 bit to indicate the presence of a trailing 32-bit checksum,
	//   * K: 10 bits to represent len(key).

	var (
		mxck uint16
	)

	e = binary.Read(decoder.reader, binary.BigEndian, &mxck)
	if e != nil {
		return
	}

	m = int(mxck>>offsetM) * common.WordSize

	x = int(mxck>>offsetX) & 0b11

	if x == 0 {
		x = 4
	}

	c = (mxck>>offsetC)&1 == 1

	k = int(mxck) & 0b1111111111

	return
}

func (decoder *Decoder) readV(x int) (v int, e error) {
	// Reads x bytes and returns the interpreted len(value).

	_, e = io.ReadFull(decoder.reader,
		decoder.buffer[maxUintLen32-x:],
	)
	if e != nil {
		return
	}

	v = int(
		binary.BigEndian.Uint32(decoder.buffer),
	)

	switch x {
	case 1:
		v &= 0x000000ff

	case 2:
		v &= 0x0000ffff

	case 3:
		v &= 0x00ffffff
	}

	return
}

func (decoder *Decoder) verifyChecksum(key, metadata, value []byte) (e error) {
	// Reads and verifies a 32-bit checksum of the record if decoder.hasher is
	// not nil; discards four bytes otherwise.

	var (
		computed uint32
		observed uint32
	)

	if decoder.hasher == nil {
		_, e = io.CopyN(io.Discard, decoder.reader, maxUintLen32)

		return
	}

	e = binary.Read(decoder.reader, binary.BigEndian, &observed)
	if e != nil {
		return
	}

	defer decoder.hasher.Reset()

	_, e = decoder.hasher.Write(key)
	if e != nil {
		return
	}

	_, e = decoder.hasher.Write(metadata)
	if e != nil {
		return
	}

	_, e = decoder.hasher.Write(value)
	if e != nil {
		return
	}

	computed = decoder.hasher.Sum32()

	if computed != observed {
		return syscall.EBADMSG
	}

	return
}
