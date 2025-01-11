package voidDB

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/voidDB/voidDB/common"
)

const (
	pageSize    = common.PageSize
	pointerSize = common.PointerSize

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

func (meta *Meta) makeCopy() (copi Meta) {
	copi = make([]byte, pageSize)

	copy(copi, *meta)

	return
}

func (meta *Meta) _magic() []byte {
	const (
		size = 8
	)

	return (*meta)[:size]
}

func (meta *Meta) setMagic() {
	copy(
		meta._magic(),
		[]byte(metaMagic),
	)

	return
}

func (meta *Meta) _version() []byte {
	const (
		offset = 8
		size   = 8
	)

	return (*meta)[offset : offset+size]
}

func (meta *Meta) version() int {
	return int(
		binary.BigEndian.Uint64(
			meta._version(),
		),
	)
}

func (meta *Meta) setVersion() {
	binary.BigEndian.PutUint64(
		meta._version(),
		uint64(version),
	)

	return
}

func (meta *Meta) isMeta() bool {
	return bytes.Equal(
		meta._magic(),
		[]byte(metaMagic),
	) &&
		meta.version() == version
}

func (meta *Meta) _timestamp() []byte {
	const (
		offset = 16
		size   = 8
	)

	return (*meta)[offset : offset+size]
}

func (meta *Meta) timestamp() time.Time {
	return time.Unix(0,
		int64(
			binary.BigEndian.Uint64(
				meta._timestamp(),
			),
		),
	)
}

func (meta *Meta) setTimestamp() {
	binary.BigEndian.PutUint64(
		meta._timestamp(),
		uint64(time.Now().UnixNano()),
	)

	return
}

func (meta *Meta) _serialNumber() []byte {
	const (
		offset = 24
		size   = 8
	)

	return (*meta)[offset : offset+size]
}

func (meta *Meta) serialNumber() int {
	return int(
		binary.BigEndian.Uint64(
			meta._serialNumber(),
		),
	)
}

func (meta *Meta) setSerialNumber(number int) {
	binary.BigEndian.PutUint64(
		meta._serialNumber(),
		uint64(number),
	)

	return
}

func (meta *Meta) _rootNodePointer() []byte {
	const (
		offset = 32
	)

	return (*meta)[offset : offset+pointerSize]
}

func (meta *Meta) rootNodePointer() int {
	return int(
		binary.BigEndian.Uint64(
			meta._rootNodePointer(),
		),
	)
}

func (meta *Meta) setRootNodePointer(pointer int) {
	binary.BigEndian.PutUint64(
		meta._rootNodePointer(),
		uint64(pointer),
	)

	return
}

func (meta *Meta) _frontierPointer() []byte {
	const (
		offset = 40
	)

	return (*meta)[offset : offset+pointerSize]
}

func (meta *Meta) frontierPointer() int {
	return int(
		binary.BigEndian.Uint64(
			meta._frontierPointer(),
		),
	)
}

func (meta *Meta) setFrontierPointer(size int) {
	binary.BigEndian.PutUint64(
		meta._frontierPointer(),
		uint64(size),
	)

	return
}
