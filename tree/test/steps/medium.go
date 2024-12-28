package steps

import (
	//"encoding/hex"
	//"log"

	"github.com/voidDB/voidDB/tree"
)

type Medium []byte

func (m *Medium) Save(bytes []byte) (pointer int, e error) {
	pointer = len(*m)

	*m = append(*m, bytes...)

	*m = append(*m, // padding
		make([]byte,
			tree.PageSize-(len(bytes)%tree.PageSize),
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
