package cursor

type Medium interface {
	Save([]byte) int
	Load(int, int) []byte
	Free(int, int)
}
