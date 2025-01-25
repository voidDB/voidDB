package voidDB

import (
	"bytes"
	"time"

	"github.com/voidDB/voidDB/common"
)

const (
	pageSize = common.PageSize
	lineSize = common.LineSize
	wordSize = common.WordSize

	metaMagic = "voidMETA"
	version   = 0
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
	return bytes.Equal(
		meta.magic(),
		[]byte(metaMagic),
	) &&
		meta.getVersion() == version
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

func (meta voidMeta) freeQueue(size int) freeQueue {
	return common.Field(meta,
		pageSize/2+lineSize*(logarithm(size)-1),
		lineSize,
	)
}

type freeQueue []byte

func (queue freeQueue) headPointer() []byte {
	return common.Field(queue, 0, wordSize)
}

func (queue freeQueue) getHeadPointer() int {
	return common.GetInt(
		queue.headPointer(),
	)
}

func (queue freeQueue) setHeadPointer(pointer int) {
	common.PutInt(
		queue.headPointer(),
		pointer,
	)

	return
}

func (queue freeQueue) nextIndex() []byte {
	return common.Field(queue, wordSize, wordSize)
}

func (queue freeQueue) getNextIndex() int {
	return common.GetInt(
		queue.nextIndex(),
	)
}

func (queue freeQueue) setNextIndex(pointer int) {
	common.PutInt(
		queue.nextIndex(),
		pointer,
	)

	return
}

func (queue freeQueue) tailPointer() []byte {
	return common.Field(queue, 2*wordSize, wordSize)
}

func (queue freeQueue) getTailPointer() int {
	return common.GetInt(
		queue.tailPointer(),
	)
}

func (queue freeQueue) setTailPointer(pointer int) {
	common.PutInt(
		queue.tailPointer(),
		pointer,
	)

	return
}
