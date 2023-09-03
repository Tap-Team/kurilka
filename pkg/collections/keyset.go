package collections

import (
	"slices"
)

type KeySet[T Keyable] struct {
	// map key to index element in collection
	keys     map[int]int
	elements []T
}

func NewKeySet[T Keyable](elements []T) *KeySet[T] {
	set := &KeySet[T]{
		keys:     map[int]int{},
		elements: make([]T, 0, len(elements)),
	}
	for i := range elements {
		set.Add(elements[i])
	}
	return set
}

func (s *KeySet[T]) Add(element T) bool {
	if s.Exist(element) {
		return false
	}
	s.elements = append(s.elements, element)
	s.keys[element.Key()] = len(s.elements) - 1
	return true
}

func (s *KeySet[T]) RemoveIndex(index int) bool {
	if len(s.elements)-1 < index {
		return false
	}
	key := s.elements[index].Key()
	delete(s.keys, key)
	s.elements = slices.Delete(s.elements, index, index+1)
	return true
}

func (s *KeySet[T]) Remove(element T) bool {
	if !s.Exist(element) {
		return false
	}
	index := s.keys[element.Key()]
	s.elements = slices.Delete(s.elements, index, index+1)
	delete(s.keys, element.Key())
	for i := index + 1; i < len(s.elements); i++ {
		key := s.elements[i].Key()
		if _, ok := s.keys[key]; ok {
			s.keys[key]--
		}
	}
	return true
}

func (s *KeySet[T]) Set(index int, element T) bool {
	if s.Exist(element) {
		return false
	}
	s.elements = append(s.elements, element)
	s.keys[element.Key()] = len(s.elements) - 1
	return true
}

func (s *KeySet[T]) Exist(element T) bool {
	_, ok := s.keys[element.Key()]
	return ok
}

func (s *KeySet[T]) Element(i int) T {
	return s.elements[i]
}

func (s *KeySet[T]) Index(element T) int {
	return s.keys[element.Key()]
}

func (s *KeySet[T]) Len() int {
	return len(s.elements)
}

func (s *KeySet[T]) Cap() int {
	return cap(s.elements)
}

func (s *KeySet[T]) Iterator() Iterator[T] {
	return NewIterator[T](s.elements)
}

func (s *KeySet[T]) Elements() []T {
	elements := make([]T, len(s.elements))
	copy(elements, s.elements)
	return elements
}
