package cursor

type Medium interface {
    //Make() []byte // TODO
	Meta() []byte
	Save([]byte) int
	Load(int, int) []byte
	Free(int, int)
}
