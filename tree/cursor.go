package tree

type Cursor struct {
	medium Medium

	offset int
	index  int

	stack []ancestor
}
