package voidDB

import (
	"bytes"
	"hash"
	"hash/fnv"
	"time"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/fifo"
)

const (
	pageSize = common.PageSize
	lineSize = common.LineSize
	wordSize = common.WordSize

	metaMagic = "voidMETA"
	version   = 0
)

var (
	fnvHash hash.Hash64 = fnv.New64a()
)

type voidMeta []byte

func newMeta() (meta voidMeta) {
	meta = make([]byte, pageSize)

	meta.setMagic()

	meta.setVersion()

	return
}

func newMetaInit() (meta voidMeta) {
	meta = newMeta()

	meta.setTimestamp()

	meta.setRootNodePointer(2 * pageSize)

	meta.setFrontierPointer(3 * pageSize)

	meta.setChecksum()

	return
}

func (meta voidMeta) makeCopy() (copi voidMeta) {
	copi = make([]byte, pageSize)

	copy(copi, meta)

	return
}

func (meta voidMeta) magic() []byte {
	return common.Field(meta, 0, wordSize)
}

func (meta voidMeta) setMagic() {
	copy(
		meta.magic(),
		[]byte(metaMagic),
	)

	return
}

func (meta voidMeta) version() []byte {
	return common.Field(meta, wordSize, wordSize)
}

func (meta voidMeta) getVersion() int {
	return common.GetInt(
		meta.version(),
	)
}

func (meta voidMeta) setVersion() {
	common.PutInt(
		meta.version(),
		version,
	)

	return
}

func (meta voidMeta) isMeta() bool {
	switch {
	case !meta.checksumOK():
		return false

	case !bytes.Equal(meta.magic(), []byte(metaMagic)):
		return false

	case meta.getVersion() != version:
		return false
	}

	return true
}

func (meta voidMeta) timestamp() []byte {
	return common.Field(meta, 2*wordSize, wordSize)
}

func (meta voidMeta) getTimestamp() time.Time {
	return time.Unix(0,
		int64(
			common.GetInt(
				meta.timestamp(),
			),
		),
	)
}

func (meta voidMeta) setTimestamp() {
	common.PutInt(
		meta.timestamp(),
		int(time.Now().UnixNano()),
	)

	return
}

func (meta voidMeta) serialNumber() []byte {
	return common.Field(meta, 3*wordSize, wordSize)
}

func (meta voidMeta) getSerialNumber() int {
	return common.GetInt(
		meta.serialNumber(),
	)
}

func (meta voidMeta) setSerialNumber(number int) {
	common.PutInt(
		meta.serialNumber(),
		number,
	)

	return
}

func (meta voidMeta) rootNodePointer() []byte {
	return common.Field(meta, 4*wordSize, wordSize)
}

func (meta voidMeta) getRootNodePointer() int {
	return common.GetInt(
		meta.rootNodePointer(),
	)
}

func (meta voidMeta) setRootNodePointer(pointer int) {
	common.PutInt(
		meta.rootNodePointer(),
		pointer,
	)

	return
}

func (meta voidMeta) frontierPointer() []byte {
	return common.Field(meta, 5*wordSize, wordSize)
}

func (meta voidMeta) getFrontierPointer() int {
	return common.GetInt(
		meta.frontierPointer(),
	)
}

func (meta voidMeta) setFrontierPointer(pointer int) {
	common.PutInt(
		meta.frontierPointer(),
		pointer,
	)

	return
}

func (meta voidMeta) checksum() []byte {
	return common.Field(meta, 7*wordSize, wordSize)
}

func (meta voidMeta) computeChecksum() []byte {
	fnvHash.Reset()

	fnvHash.Write(meta[:7*wordSize])

	fnvHash.Write(
		make([]byte, wordSize),
	)

	fnvHash.Write(meta[8*wordSize:])

	return fnvHash.Sum(nil)
}

func (meta voidMeta) setChecksum() {
	copy(
		meta.checksum(),
		meta.computeChecksum(),
	)

	return
}

func (meta voidMeta) checksumOK() bool {
	return bytes.Equal(
		meta.checksum(),
		meta.computeChecksum(),
	)
}

func (meta voidMeta) freeQueue(size int) fifo.FIFO {
	return common.Field(meta,
		pageSize/2+lineSize*(logarithm(size)-1),
		lineSize,
	)
}
