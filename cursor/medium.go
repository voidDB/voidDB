package cursor

type Medium interface {
	Meta() []byte
	Save([]byte) int
	Load(int, int) []byte
	Free(int, int)
}
