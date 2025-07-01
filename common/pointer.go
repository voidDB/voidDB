package common

const (
	// HACK: Pointers to data are multiples of wordSize, hence the least
	// significant three bits are free to mean other things.
	Tombstone = 1
	graveyard = 2
	undefined = 4 // for future use
	flagsMask = 7
)

type Pointer int

func (pointer Pointer) mask(m int) int {
	return int(pointer) &^ m
}

func (pointer Pointer) Clean() int {
	return pointer.mask(flagsMask)
}

func (pointer Pointer) hasFlag(flag int) bool {
	return int(pointer)&flag > 0
}

func (pointer Pointer) IsTombstone() bool {
	return pointer == Tombstone
}

func (pointer Pointer) IsGraveyard() bool {
	return pointer.hasFlag(graveyard)
}

func (pointer Pointer) IsDeleted() bool {
	return pointer.hasFlag(Tombstone | graveyard)
}

func (pointer Pointer) setFlag(flag int) int {
	return int(pointer) | flag
}

func (pointer Pointer) ToGraveyard() int {
	return pointer.setFlag(graveyard)
}
