package cursor

type Medium interface {
	Meta() []byte
	Save([]byte) int
	Page(int) ([]byte, bool)
	Data(int, int) []byte
	Free(int, int)
	Root(int)
}
