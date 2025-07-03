package fifo

type Medium interface {
	Save([]byte) int
	SaveAt(int, []byte)
	Page(int) ([]byte, bool)
	Free(int, int)
}
