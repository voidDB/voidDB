package tree

const (
	tombstone = 1 // Pointers are multiples of PageSize,
	// hence the least significant 11 bits are free to mean other things.
)

func Del(medium Medium, offset int, key []byte) (pointer int, e error) {
	var (
		node Node = getNode(medium, offset, true)

		index  int
		valLen int
	)

	index, pointer, valLen = node.search(key)

	switch {
	case pointer == tombstone:
		fallthrough

	case pointer == 0:
		return 0, ErrorNotFound

	case valLen > 0:
		node = node.update(index, tombstone, 0)

		medium.Free(pointer, valLen)

	default:
		pointer, e = Del(medium, pointer, key)
		if e != nil {
			return
		}

		node = node.update(index, pointer, 0)
	}

	return medium.Save(node)
}
