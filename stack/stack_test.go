package stack_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/stack"
)

func TestPopFromEmptyStackReturnsFalse(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	_, got := stk.Pop()
	want := false
	if want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}

func TestPopFromNonEmptyStackReturnsTopItem(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	stk.Push(1)
	stk.Push(2)
	got, ok := stk.Pop()
	if !ok {
		t.Fatal("expected pop ok status to be true")
	}
	want := 2
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSizeOfEmptyStackReturnsExpectedSize(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[string]
	got := stk.Size()
	want := 0
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSizeOfNonEmptyStackReturnsExpectedSize(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[string]
	stk.Push("a")
	stk.Push("b")
	stk.Push("c")
	got := stk.Size()
	want := 3
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSizeChangesAfterItemsPushedOntoStack(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[bool]
	initSize := stk.Size()
	stk.Push(true)
	finalSize := stk.Size()
	want := 1
	got := finalSize - initSize
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSizeChangesAfterItemsPoppedFromStack(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	for i := 0; i < 4; i++ {
		stk.Push(i)
	}
	initSize := stk.Size()
	for i := 0; i < 2; i++ {
		_, ok := stk.Pop()
		if !ok {
			t.Fatal("did not expect stack pop status to be false")
		}
	}
	finalSize := stk.Size()
	got := initSize - finalSize
	want := 2
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestPeekFromEmptyStackReturnsFalse(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	_, got := stk.Peek()
	want := false
	if want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}

func TestPeekFromNonEmptyStackReturnsTopItem(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	stk.Push(1)
	stk.Push(2)
	got, ok := stk.Peek()
	if !ok {
		t.Fatal("expected peek ok status to be true")
	}
	want := 2
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestPeekFromNonEmptyStackDoesNotChangeStackSize(t *testing.T) {
	t.Parallel()
	var stk stack.Stack[int]
	stk.Push(1)
	initSize := stk.Size()
	_, ok := stk.Peek()
	if !ok {
		t.Fatal("expected peek ok status to be true")
	}
	finalSize := stk.Size()
	if initSize != finalSize {
		t.Errorf("size before peek (%d) does not equal size after peek (%d)", initSize, finalSize)
	}
}
