package stack

type Stack[T any] struct {
	arr  []T
	size int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Empty() bool {
	return s.size == 0
}

func (s *Stack[T]) Top() T {
	return s.arr[s.size-1]
}

func (s *Stack[T]) Push(val T) {
	if s.size < cap(s.arr) {
		s.arr = append(s.arr, val)
		s.arr[s.size] = val
	} else {
		s.arr = append(s.arr, val)
	}
	s.size++
}

func (s *Stack[T]) Pop() T {
	s.size--
	return s.arr[s.size]
}
