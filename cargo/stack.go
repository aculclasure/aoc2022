package cargo

import (
	"sync"
)

// Stack represents a concurrency-safe stack for pushing cargo items onto and
// popping cargo items off from.
type Stack struct {
	mtx   *sync.Mutex
	items []rune
}

// NewStack accepts an optional number of initial items, pushes them onto a stack
// and returns the stack. An empty stack is returned if no initial items are given.
func NewStack(items ...rune) *Stack {
	stk := &Stack{mtx: &sync.Mutex{}}
	for _, v := range items {
		stk.Push(v)
	}
	return stk
}

// Push accepts an item and pushes it onto the stack.
func (s *Stack) Push(item rune) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.items = append(s.items, item)
}

// Pop removes the top item from a stack and returns it along with a boolean
// value indicating if there was an item to pop. The boolean value is false when
// a Pop is attempted on an empty stack.
func (s *Stack) Pop() (rune, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if len(s.items) == 0 {
		return '0', false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

// Peek returns the top item of the stack without removing it from the stack
// along with a boolean value indicating if the stack is empty or not. The
// boolean value is false when the stack is empty.
func (s *Stack) Peek() (rune, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if len(s.items) == 0 {
		return '0', false
	}

	index := len(s.items) - 1
	item := s.items[index]
	return item, true
}

// Size returns the number of items in the stack.
func (s *Stack) Size() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return len(s.items)
}

// Items returns a copy slice of the items stored in the stack. The beginning of
// the returned slice represents the bottom of the stack.
func (s *Stack) Items() []rune {
	var stkItems []rune
	s.mtx.Lock()
	defer s.mtx.Unlock()
	stkItems = append(stkItems, s.items...)
	return stkItems
}
