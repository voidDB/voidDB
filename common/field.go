package common

func Field(slice []byte, offset, length int) []byte {
	return slice[offset : offset+length]
}
