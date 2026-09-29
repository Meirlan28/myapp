package optional

type Optional[T any] struct {
	Value *T
	Set   bool
}
