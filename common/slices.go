package common

const (
	PageSize = 4096
	KeySize  = 512
	LineSize = 64
	TwinSize = 16
	WordSize = 8
	HalfSize = 4
)

func NewPage() []byte {
	return make([]byte, PageSize)
}

func NewWord() []byte {
	return make([]byte, WordSize)
}

func Slice(super []byte, offset, length int) []byte {
	return super[offset : offset+length]
}

func Page(super []byte, offset int) []byte {
	return Slice(super, offset, PageSize)
}

func PageN(super []byte, index int) []byte {
	return Page(super, index*PageSize)
}

func Key(super []byte, offset int) []byte {
	return Slice(super, offset, KeySize)
}

func KeyN(super []byte, index int) []byte {
	return Key(super, index*KeySize)
}

func Line(super []byte, offset int) []byte {
	return Slice(super, offset, LineSize)
}

func LineN(super []byte, index int) []byte {
	return Line(super, index*LineSize)
}

func Twin(super []byte, offset int) []byte {
	return Slice(super, offset, TwinSize)
}

func TwinN(super []byte, index int) []byte {
	return Twin(super, index*TwinSize)
}

func Word(super []byte, offset int) []byte {
	return Slice(super, offset, WordSize)
}

func WordN(super []byte, index int) []byte {
	return Word(super, index*WordSize)
}

func Half(super []byte, offset int) []byte {
	return Slice(super, offset, HalfSize)
}

func HalfN(super []byte, index int) []byte {
	return Half(super, index*HalfSize)
}
