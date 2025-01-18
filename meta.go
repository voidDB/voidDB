package voidDB

import (
	"bytes"
	"time"

	"github.com/voidDB/voidDB/common"
)

const (
	pageSize = common.PageSize
	wordSize = common.WordSize

	metaMagic = "voidMETA"
	version   = 0
)

type Meta []byte

func newMeta() (meta Meta) {
	meta = make([]byte, pageSize)

	meta.setMagic()

	meta.setVersion()

	return
}

func newMetaInit() (meta Meta) {
	meta = newMeta()

	meta.setTimestamp()

	meta.setRootNodePointer(pageSize)

	meta.setFrontierPointer(2 * pageSize)

	return
}

func (meta Meta) makeCopy() (copi Meta) {
	copi = make([]byte, pageSize)

	copy(copi, meta)

	return
}

func (meta Meta) magic() []byte {
	return common.Field(meta, 0, wordSize)
}

func (meta Meta) setMagic() {
	copy(
		meta.magic(),
		[]byte(metaMagic),
	)

	return
}

func (meta Meta) version() []byte {
	return common.Field(meta, wordSize, wordSize)
}

func (meta Meta) getVersion() int {
	return common.GetInt(
		meta.version(),
	)
}

func (meta Meta) setVersion() {
	common.PutInt(
		meta.version(),
		version,
	)

	return
}

func (meta Meta) isMeta() bool {
	return bytes.Equal(
		meta.magic(),
		[]byte(metaMagic),
	) &&
		meta.getVersion() == version
}

func (meta Meta) timestamp() []byte {
	return common.Field(meta, 2*wordSize, wordSize)
}

func (meta Meta) getTimestamp() time.Time {
	return time.Unix(0,
		int64(
			common.GetInt(
				meta.timestamp(),
			),
		),
	)
}

func (meta Meta) setTimestamp() {
	common.PutInt(
		meta.timestamp(),
		int(time.Now().UnixNano()),
	)

	return
}

func (meta Meta) serialNumber() []byte {
	return common.Field(meta, 3*wordSize, wordSize)
}

func (meta Meta) getSerialNumber() int {
	return common.GetInt(
		meta.serialNumber(),
	)
}

func (meta Meta) setSerialNumber(number int) {
	common.PutInt(
		meta.serialNumber(),
		number,
	)

	return
}

func (meta Meta) rootNodePointer() []byte {
	return common.Field(meta, 4*wordSize, wordSize)
}

func (meta Meta) getRootNodePointer() int {
	return common.GetInt(
		meta.rootNodePointer(),
	)
}

func (meta Meta) setRootNodePointer(pointer int) {
	common.PutInt(
		meta.rootNodePointer(),
		pointer,
	)

	return
}

func (meta Meta) frontierPointer() []byte {
	return common.Field(meta, 5*wordSize, wordSize)
}

func (meta Meta) getFrontierPointer() int {
	return common.GetInt(
		meta.frontierPointer(),
	)
}

func (meta Meta) setFrontierPointer(size int) {
	common.PutInt(
		meta.frontierPointer(),
		size,
	)

	return
}