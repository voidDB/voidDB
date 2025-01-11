package steps

import (
	//"encoding/hex"
	//"log"

	"github.com/voidDB/voidDB/common"
)

type Medium []byte

func (m *Medium) Save(bytes []byte) (pointer int) {
	pointer = len(*m)

	*m = append(*m, bytes...)

	*m = append(*m, // padding
		make([]byte,
			common.PageSize-(len(bytes)%common.PageSize),
		)...,
	)

	//log.Println(
	//	hex.Dump(*m),
	//)

	return
}

func (m *Medium) Load(offset, length int) []byte {
	return (*m)[offset : offset+length]
}

func (m *Medium) Free(offset, length int) {
	return
}
