package free

func Put(medium Medium, offset, txnID int, pointers []int) (head, tail int) {
	var (
		i      int
		length int

		free Free = NewFree(txnID)
	)

	switch {
	case len(pointers) == 0:
		tail = medium.Save(free)

		return tail, tail

	case len(pointers) > MaxNodeLength:
		length = MaxNodeLength

	default:
		length = len(pointers)
	}

	free.setLength(length)

	for i = 0; i < length; i++ {
		free.setPagePointer(i,
			pointers[i],
		)
	}

	head, tail = Put(medium, -1, txnID, pointers[length:])

	free.setNextPointer(head)

	if offset > 0 {
		medium.SaveAt(offset, free)

		return offset, tail
	}

	return medium.Save(free), tail
}
