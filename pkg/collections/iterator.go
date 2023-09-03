package collections

type Iterator[T any] interface {
	Next() bool
	Value() T
	Index() int
}

type iterator[T any] struct {
	elements []T
	index    int
}

func (i *iterator[T]) Next() bool {
	i.index++
	if len(i.elements) <= i.index {
		return false
	}
	return true
}

func (i *iterator[T]) Value() T {
	return i.elements[i.index]
}

func (i *iterator[T]) Index() int {
	return i.index
}

func NewIterator[T any](elements []T) Iterator[T] {
	return &iterator[T]{elements: elements, index: -1}
}
