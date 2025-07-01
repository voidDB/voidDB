package node

import (
	"github.com/voidDB/voidDB/common"
)

var (
	nodeMagic = []byte("voidNODE")
)

type Meta []byte

func (meta Meta) magic() []byte {
	return common.WordN(meta, 0)
}

func (meta Meta) setMagic() {
	copy(
		meta.magic(),
		nodeMagic,
	)

	return
}

func (meta Meta) vetMagic() error {
	return common.ErrorIfNotEqual(
		meta.magic(),
		nodeMagic,
		common.ErrorCorrupt,
	)
}

func (meta Meta) length() []byte {
	return common.WordN(meta, 1)
}

func (meta Meta) getLength() int {
	return common.GetIntFromWord(
		meta.length(),
	)
}

func (meta Meta) setLength(length int) {
	common.PutIntIntoWord(
		meta.length(),
		length,
	)

	return
}
