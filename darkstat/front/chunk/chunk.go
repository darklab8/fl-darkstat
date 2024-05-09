package chunk

type Chunk[T any] struct {
	Items      []T
	CurrentUrl string
}
