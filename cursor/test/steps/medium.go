package steps

import (
	//"encoding/hex"
	//"log"

	"github.com/voidDB/voidDB/common"
)

type Medium []byte

func (m *Medium) Meta() []byte {
	return []byte("voidTestMetadata")
}

func (m *Medium) Save(bytes []byte) (pointer int) {
	pointer = len(*m)

	*m = append(*m, bytes...)

	if len(bytes)%common.PageSize > 0 {
		*m = append(*m, // padding
			make([]byte,
				common.PageSize-(len(bytes)%common.PageSize),
			)...,
		)
	}

	//log.Println(
	//	hex.Dump(*m),
	//)

	return
}

func (m *Medium) Page(offset int) ([]byte, bool) {
	return common.Page(*m, offset), false
}

func (m *Medium) Data(offset, length int) []byte {
	return common.Slice(*m, offset, length)
}

func (m *Medium) Free(offset, length int) {
	return
}

func (m *Medium) Root(int) {
	return
}
