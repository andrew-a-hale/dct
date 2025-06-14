package utils

type Stack[T any] []T

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) Pop() T {
	i := s.Len() - 1
	x := (*s)[i]
	*s = (*s)[:i]
	return x
}

func (s *Stack[any]) Push(el ...any) {
	*s = append(*s, el...)
}

func (s *Stack[T]) Peek() T {
	return (*s)[s.Len()-1]
}
