package fifo

func (fifo FIFO) Enqueue(medium Medium, txnID int, pointers []int,
	cannibalise bool,
) {
	var (
		enFn enqueueFunc
		head int
		tail int
	)

	switch {
	case cannibalise:
		enFn = enqueueCannibalise

	default:
		enFn = enqueue
	}

	head, tail = enFn(medium,
		fifo.getTailPointer(),
		txnID,
		pointers,
	)

	if fifo.getHeadPointer() == 0 {
		fifo.setHeadPointer(head)
	}

	fifo.setTailPointer(tail)

	return
}

type enqueueFunc func(Medium, int, int, []int) (int, int)

func enqueue(medium Medium, offset, txnID int, pointers []int) (
	head, tail int,
) {
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
		free.setPagePointer(i, pointers[i])
	}

	head, tail = enqueue(medium, -1, txnID, pointers[length:])

	free.setNextPointer(head)

	if offset > 0 {
		medium.SaveAt(offset, free)

		return offset, tail
	}

	return medium.Save(free), tail
}

func enqueueCannibalise(medium Medium, offset, txnID int, pointers []int) (
	head, tail int,
) {
	var (
		i      int
		length int

		free Free = NewFree(txnID)
	)

	switch {
	case len(pointers) == 0:
		tail = offset

		goto end

	case len(pointers) > MaxNodeLength:
		length = MaxNodeLength

	default:
		length = len(pointers) - 1
	}

	free.setLength(length)

	for i = 0; i < length; i++ {
		free.setPagePointer(i, pointers[i])
	}

	head, tail = enqueueCannibalise(medium, pointers[length], txnID,
		pointers[length+1:],
	)

	free.setNextPointer(head)

end:
	medium.SaveAt(offset, free)

	return offset, tail
}
