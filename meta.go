package voidDB

import (
	"bytes"
	"time"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/fifo"
)

const (
	version = 1 // broken compatibility with v0.1.x
)

var (
	metaMagic = []byte("voidMETA")
)

type voidMeta []byte

func newMeta() (meta voidMeta) {
	meta = common.NewPage()

	meta.setMagic()

	meta.setVersion()

	return
}

func newMetaInit() (meta voidMeta) {
	meta = newMeta()

	meta.setTimestamp()

	meta.setRootNodePointer(2 * common.PageSize)

	meta.setFrontierPointer(3 * common.PageSize)

	return
}

func (meta voidMeta) makeCopy(concise bool) (copi voidMeta) {
	switch concise {
	case true:
		copi = common.NewLine()

	default:
		copi = common.NewPage()
	}

	copy(copi, meta)

	return
}

func (meta voidMeta) magic() []byte {
	return common.WordN(meta, 0)
}

func (meta voidMeta) setMagic() {
	copy(
		meta.magic(),
		metaMagic,
	)

	return
}

func (meta voidMeta) version() []byte {
	return common.WordN(meta, 1)
}

func (meta voidMeta) getVersion() int {
	return common.GetIntFromWord(
		meta.version(),
	)
}

func (meta voidMeta) setVersion() {
	common.PutIntIntoWord(
		meta.version(),
		version,
	)

	return
}

func (meta voidMeta) isMeta() bool {
	switch {
	case !bytes.Equal(meta.magic(), metaMagic):
		return false

	case meta.getVersion() != version:
		return false
	}

	return true
}

func (meta voidMeta) timestamp() []byte {
	return common.WordN(meta, 2)
}

func (meta voidMeta) getTimestamp() time.Time {
	return time.Unix(0,
		int64(
			common.GetIntFromWord(
				meta.timestamp(),
			),
		),
	)
}

func (meta voidMeta) setTimestamp() {
	common.PutIntIntoWord(
		meta.timestamp(),
		int(time.Now().UnixNano()),
	)

	return
}

func (meta voidMeta) serialNumber() []byte {
	return common.WordN(meta, 3)
}

func (meta voidMeta) getSerialNumber() int {
	return common.GetIntFromWord(
		meta.serialNumber(),
	)
}

func (meta voidMeta) setSerialNumber(number int) {
	common.PutIntIntoWord(
		meta.serialNumber(),
		number,
	)

	return
}

func (meta voidMeta) rootNodePointer() []byte {
	return common.WordN(meta, 4)
}

func (meta voidMeta) getRootNodePointer() int {
	return common.GetIntFromWord(
		meta.rootNodePointer(),
	)
}

func (meta voidMeta) setRootNodePointer(pointer int) {
	common.PutIntIntoWord(
		meta.rootNodePointer(),
		pointer,
	)

	return
}

func (meta voidMeta) frontierPointer() []byte {
	return common.WordN(meta, 5)
}

func (meta voidMeta) getFrontierPointer() int {
	return common.GetIntFromWord(
		meta.frontierPointer(),
	)
}

func (meta voidMeta) setFrontierPointer(pointer int) {
	common.PutIntIntoWord(
		meta.frontierPointer(),
		pointer,
	)

	return
}

func (meta voidMeta) freeQueue(size int) fifo.FIFO {
	var (
		queues = common.Slice(meta, common.PageSize/2, common.PageSize/2)
	)

	return common.LineN(queues, logarithm(size)-1)
}
