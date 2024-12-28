package tree

type Medium interface {
	Save([]byte) (int, error)
	Load(int, int) []byte
}
