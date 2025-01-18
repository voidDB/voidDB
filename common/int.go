package common

import (
	"encoding/binary"
)

func GetInt(slice []byte) int {
	switch len(slice) {
	case WordSize:
		return int(
			binary.BigEndian.Uint64(slice),
		)

	case HalfSize:
		return int(
			binary.BigEndian.Uint32(slice),
		)
	}

	return -1
}

func PutInt(slice []byte, i int) {
	switch len(slice) {
	case WordSize:
		binary.BigEndian.PutUint64(slice,
			uint64(i),
		)

	case HalfSize:
		binary.BigEndian.PutUint32(slice,
			uint32(i),
		)
	}

	return
}
