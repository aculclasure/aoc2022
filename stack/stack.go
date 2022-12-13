// Package stack implements a generic stack data structure.
package stack

// Stack represents a generic stack data structure.
type Stack[T any] struct {
	vals []T
}

// Push accepts a value T and pushes it onto the stack.
func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

// Pop removes and returns the top item from the stack along with a boolean
// value indicating if the Pop was successful. The boolean value will be false
// when attempting to pop from an empty stack.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}

	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}

	return s.vals[len(s.vals)-1], true
}

// Size returns the number of items in the stack.
func (s Stack[T]) Size() int {
	return len(s.vals)
}
