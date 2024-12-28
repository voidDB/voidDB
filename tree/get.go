package tree

func Get(medium Medium, offset int, key []byte) (value []byte, e error) {
	var (
		node Node = getNode(medium, offset)

		pointer int
		valLen  int
	)

	_, pointer, valLen = node.search(key)

	switch {
	case pointer == 0:
		e = ErrorNotFound

		return

	case valLen > 0:
		value = medium.Load(pointer, valLen)

		return
	}

	return Get(medium, pointer, key)
}
