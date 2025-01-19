package free

type Medium interface {
	Save([]byte) int
	SaveAt(int, []byte)
	Load(int, int) []byte
	Free(int, int)
}
