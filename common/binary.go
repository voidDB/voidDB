package common

import (
	"encoding/binary"
)

func GetIntFromWord(word []byte) int {
	return int(
		binary.BigEndian.Uint64(word),
	)
}

func GetIntFromHalf(half []byte) int {
	return int(
		binary.BigEndian.Uint32(half),
	)
}

func PutIntIntoWord(word []byte, i int) {
	binary.BigEndian.PutUint64(word,
		uint64(i),
	)

	return
}

func PutIntIntoHalf(half []byte, i int) {
	binary.BigEndian.PutUint32(half,
		uint32(i),
	)

	return
}
