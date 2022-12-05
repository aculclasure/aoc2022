package cargo_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/cargo"
	"github.com/google/go-cmp/cmp"
)

func TestNewStackWithNoInitialItemsReturnsEmptyStack(t *testing.T) {
	t.Parallel()
	var want []rune
	stk := cargo.NewStack()
	got := stk.Items()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestNewStackWithInitialItemsReturnsStackInitializedWithItems(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('A', 'B', 'C')
	want := []rune{'A', 'B', 'C'}
	got := stk.Items()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestStack_PopReturnsTopOfStack(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('A', 'B', 'C')
	got, ok := stk.Pop()
	if !ok {
		t.Fatal("expected pop status to be true")
	}

	want := 'C'
	if want != got {
		t.Errorf("want %s, got %s", string(want), string(got))
	}
}

func TestStack_PopRemovesItemFromStack(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('A', 'B', 'C')
	_, ok := stk.Pop()
	if !ok {
		t.Fatal("expected pop status to be true")
	}

	want := []rune{'A', 'B'}
	got := stk.Items()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestStack_PopReturnsFalseForEmptyStack(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack()
	_, got := stk.Pop()
	want := false

	if want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}

func TestStack_PopChangesStackSize(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('a', 'b', 'c')
	_, ok := stk.Pop()
	if !ok {
		t.Fatal("expected pop to return a true status")
	}

	want := 2
	got := stk.Size()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestStack_PushPushesItemToTopOfStack(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('A', 'B', 'C')
	stk.Push('D')

	want := []rune{'A', 'B', 'C', 'D'}
	got := stk.Items()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestStack_PushChangesStackSize(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack()
	stk.Push('a')

	want := 1
	got := stk.Size()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestStack_PeekFromEmptyStackReturnsFalseStatus(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack()

	_, ok := stk.Peek()
	if ok {
		t.Error("want status to be false")
	}
}

func TestStack_PeekFromNonEmptyStackReturnsExpectedItem(t *testing.T) {
	t.Parallel()
	stk := cargo.NewStack('a', 'b', 'c')
	want := 'c'
	got, ok := stk.Peek()
	if !ok {
		t.Fatal("want status to be true")
	}

	if want != got {
		t.Errorf("want %s, got %s", string(want), string(got))
	}
}

func TestStack_Size(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input cargo.Stack
		want  int
	}{
		"Empty stack returns 0": {
			input: *cargo.NewStack(),
			want:  0,
		},
		"Non-empty stack returns expected size": {
			input: *cargo.NewStack('a', 'b', 'c', 'd', 'e'),
			want:  5,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Size()
			if tc.want != got {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}
