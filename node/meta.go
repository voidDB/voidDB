package node

import (
	"bytes"

	"github.com/voidDB/voidDB/common"
)

type Meta []byte

func (meta Meta) magic() []byte {
	return common.Field(meta, 0, wordSize)
}

func (meta Meta) isNode() bool {
	return bytes.Equal(
		meta.magic(),
		[]byte(nodeMagic),
	)
}

func (meta Meta) setMagic() {
	copy(
		meta.magic(),
		[]byte(nodeMagic),
	)

	return
}

func (meta Meta) length() []byte {
	return common.Field(meta, wordSize, wordSize)
}

func (meta Meta) getLength() int {
	return common.GetInt(
		meta.length(),
	)
}

func (meta Meta) setLength(length int) {
	common.PutInt(
		meta.length(),
		length,
	)

	return
}
