package cargo_test

import (
	"io"
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/cargo"
	"github.com/google/go-cmp/cmp"
)

func TestGetCrates(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  []cargo.Crate
	}{
		"Line starting with crate and skipping column between next crate returns a valid slice of crates": {
			input: "[G]     [P]",
			want: []cargo.Crate{
				{Item: 'G', Stack: 1},
				{Item: 'P', Stack: 3},
			},
		},
		"Line starting and ending with skipped columns returns a valid slice of crates": {
			input: "            [M] [S] [S]            ",
			want: []cargo.Crate{
				{Item: 'M', Stack: 4},
				{Item: 'S', Stack: 5},
				{Item: 'S', Stack: 6},
			},
		},
		"Line with no skipped columns returns a valid slice of crates": {
			input: "[C] [G] [Q]",
			want: []cargo.Crate{
				{Item: 'C', Stack: 1},
				{Item: 'G', Stack: 2},
				{Item: 'Q', Stack: 3},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := cargo.GetCrates(tc.input)
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestInitializeFromCrateRows(t *testing.T) {
	t.Parallel()
	layout, err := cargo.NewLayout(9)
	if err != nil {
		t.Fatal(err)
	}

	input := [][]cargo.Crate{
		{
			{Item: 'M', Stack: 4},
			{Item: 'S', Stack: 5},
			{Item: 'S', Stack: 6},
		},
		{
			{Item: 'M', Stack: 3},
			{Item: 'N', Stack: 4},
			{Item: 'L', Stack: 5},
			{Item: 'T', Stack: 6},
			{Item: 'Q', Stack: 7},
		},
	}

	err = layout.InitializeFromCrateRows(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	want := [][]rune{
		nil,
		nil,
		{'M'},
		{'N', 'M'},
		{'L', 'S'},
		{'T', 'S'},
		{'Q'},
		nil,
		nil,
	}
	var got [][]rune
	for _, stk := range layout.Stacks[1:] {
		got = append(got, stk.Items())
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMovementFromLine_ErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Non-numerical quantity field returns error":       "move a from 1 to 2",
		"Non-numerical src stack field returns error":      "move 2 from x to 3",
		"Non-numerical dest stack field returns error":     "move 2 from 4 to z",
		"Negative value in quantity field returns error":   "move -1 from 1 to 2",
		"Negative value in src stack field returns error":  "move 2 from -1 to 2",
		"Negative value in dest stack field returns error": "move 2 from 1 to -2",
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := cargo.MovementFromLine(tc)
			if err == nil {
				t.Fatal("expected an error but did not get one")
			}
		})
	}
}

func TestMovementFromLine(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  cargo.Movement
	}{
		"Single-digit quantity field returns a valid movement": {
			input: "move 3 from 1 to 5",
			want:  cargo.Movement{Quantity: 3, SrcStack: 1, DestStack: 5},
		},
		"Multi-digit fields returns a valid movement": {
			input: "move 100 from 10 to 100",
			want:  cargo.Movement{Quantity: 100, SrcStack: 10, DestStack: 10},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := cargo.MovementFromLine(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestLayoutFromData(t *testing.T) {
	t.Parallel()
	layoutWithSingleMovement := `[G]    
[B]    
[D]    
[C] [G]
[P] [V]
[R] [H]
 1   2

move 1 from 1 to 2
`
	layoutWithMultipleMovements := `[G]    
[B]    
[D]    
[C] [G]
[P] [V]
[R] [H]
 1   2

move 1 from 1 to 2
move 4 from 2 to 1
`
	testCases := map[string]struct {
		input io.Reader
		want  [][]rune
	}{
		"Layout data with single move updates layout": {
			input: strings.NewReader(layoutWithSingleMovement),
			want: [][]rune{
				{'R', 'P', 'C', 'D', 'B'},
				{'H', 'V', 'G', 'G'},
			},
		},
		"Layout data with multiple moves updates layout": {
			input: strings.NewReader(layoutWithMultipleMovements),
			want: [][]rune{
				{'R', 'P', 'C', 'D', 'B', 'G', 'G', 'V', 'H'},
				nil,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			layout, err := cargo.LayoutFromData(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			var got [][]rune
			for _, stk := range layout.Stacks[1:] {
				got = append(got, stk.Items())
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestLayout_AddCrateWithInvalidCrateReturnsError(t *testing.T) {
	t.Parallel()
	numStacks := 3
	layout, err := cargo.NewLayout(numStacks)
	if err != nil {
		t.Fatal(err)
	}

	invalidStackIndex := 5
	c := cargo.Crate{Item: 'a', Stack: invalidStackIndex}
	err = layout.AddCrate(c)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}

func TestLayout_AddCrateWithValidCrateUpdatesLayout(t *testing.T) {
	t.Parallel()
	layout, err := cargo.NewLayout(1)
	if err != nil {
		t.Fatal(err)
	}

	c := cargo.Crate{Item: 'a', Stack: 1}
	err = layout.AddCrate(c)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}

	want := []rune{'a'}
	got := layout.Stacks[c.Stack].Items()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestLayout_MoveErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input cargo.Movement
	}{
		"Moving from an invalid source stack index returns error": {
			input: cargo.Movement{SrcStack: 0},
		},
		"Moving to an invalid destination stack index returns error": {
			input: cargo.Movement{SrcStack: 1, DestStack: 100, Quantity: 1},
		},
		"Moving a negative quantity size returns an error": {
			input: cargo.Movement{SrcStack: 1, DestStack: 2, Quantity: -1},
		},
		"Moving a quantity of stack items greater than the source stack size returns error": {
			input: cargo.Movement{SrcStack: 1, DestStack: 2, Quantity: 100},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			layout := &cargo.Layout{
				Stacks: []*cargo.Stack{
					nil,
					cargo.NewStack('a'),
					cargo.NewStack('b'),
				},
			}
			err := layout.Move(tc.input)
			if err == nil {
				t.Fatal("expected an error but did not get one")
			}
		})
	}
}

func TestLayout_MoveWithValidMovementUpdatesLayout(t *testing.T) {
	t.Parallel()
	layout := &cargo.Layout{
		Stacks: []*cargo.Stack{
			nil,
			cargo.NewStack('a', 'b', 'c'),
			cargo.NewStack('d', 'e', 'f'),
		},
	}
	mv := cargo.Movement{SrcStack: 1, DestStack: 2, Quantity: 2}

	layout.Move(mv)
	wantDestStackItems := []rune{'d', 'e', 'f', 'c', 'b'}
	gotDestStackItems := layout.Stacks[mv.DestStack].Items()
	if !cmp.Equal(wantDestStackItems, gotDestStackItems) {
		t.Log("destination stack did not get updated as expected")
		t.Fatal(cmp.Diff(wantDestStackItems, gotDestStackItems))
	}

	wantSrcStackItems := []rune{'a'}
	gotSrcStackItems := layout.Stacks[mv.SrcStack].Items()
	if !cmp.Equal(wantSrcStackItems, gotSrcStackItems) {
		t.Log("source stack did not get updated as expected")
		t.Fatal(cmp.Diff(wantSrcStackItems, gotSrcStackItems))
	}
}

func TestLayout_GetTopItems(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input *cargo.Layout
		want  string
	}{
		"GetTopItems() from an empty layout returns an empty string": {
			input: &cargo.Layout{},
			want:  "",
		},
		"GetTopItems from a non-empty layout returns a valid string": {
			input: &cargo.Layout{
				Stacks: []*cargo.Stack{
					nil,
					cargo.NewStack('a'),
					cargo.NewStack('b', 'c'),
					cargo.NewStack('d', 'e', 'f'),
				},
			},
			want: "acf",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.input.GetTopItems()
			if tc.want != got {
				t.Errorf("want %s, got %s", tc.want, got)
			}
		})
	}
}
