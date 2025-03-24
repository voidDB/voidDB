package fifo

type Medium interface {
	Save([]byte) int
	SaveAt(int, []byte)
	Load(int, int) ([]byte, bool)
	Free(int, int)
}
