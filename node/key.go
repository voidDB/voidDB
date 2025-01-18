package node

type Key []byte

func (key Key) get(length int) []byte {
	return key[:length]
}

func (key Key) set(source []byte) (length int) {
	return copy(key, source)
}
