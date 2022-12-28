package rope_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/rope"
	"github.com/google/go-cmp/cmp"
)

func TestRope_MoveHead(t *testing.T) {
	t.Parallel()
	type move struct {
		numRows int
		numCols int
	}
	testCases := map[string]struct {
		input []move
		want  *rope.Rope
	}{
		"Moving head end up a column moves head and tail ends as expected": {
			input: []move{{numRows: 2}},
			want: &rope.Rope{
				Head: &rope.RopeEnd{
					Row:     2,
					Col:     0,
					Visited: map[string]int{"0,0": 1, "1,0": 1, "2,0": 1},
				},
				Tail: &rope.RopeEnd{
					Row:     1,
					Col:     0,
					Visited: map[string]int{"0,0": 1, "1,0": 1},
				},
			},
		},
		"Moving head across a row moves head and tail ends as expected": {
			input: []move{{numCols: -2}},
			want: &rope.Rope{
				Head: &rope.RopeEnd{
					Row:     0,
					Col:     -2,
					Visited: map[string]int{"0,0": 1, "0,-1": 1, "0,-2": 1},
				},
				Tail: &rope.RopeEnd{
					Row:     0,
					Col:     -1,
					Visited: map[string]int{"0,0": 1, "0,-1": 1},
				},
			},
		},
		"Head revisiting positions updates head's map of visited positions as expected": {
			input: []move{{numRows: 2}, {numRows: -1}},
			want: &rope.Rope{
				Head: &rope.RopeEnd{
					Row:     1,
					Col:     0,
					Visited: map[string]int{"0,0": 1, "1,0": 2, "2,0": 1},
				},
				Tail: &rope.RopeEnd{
					Row:     1,
					Col:     0,
					Visited: map[string]int{"0,0": 1, "1,0": 1},
				},
			},
		},
		"Head moving up column and across rows moves tail as expected": {
			input: []move{{numRows: 1}, {numCols: 1}, {numCols: 1}},
			want: &rope.Rope{
				Head: &rope.RopeEnd{
					Row:     1,
					Col:     2,
					Visited: map[string]int{"0,0": 1, "1,0": 1, "1,1": 1, "1,2": 1},
				},
				Tail: &rope.RopeEnd{
					Row:     1,
					Col:     1,
					Visited: map[string]int{"0,0": 1, "1,1": 1},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rp := rope.NewRope()
			for _, mv := range tc.input {
				rp.MoveHead(mv.numRows, mv.numCols)
			}
			got := rp
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestRope_UpdateTail(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input *rope.Rope
		want  *rope.RopeEnd
	}{
		"Head and tail in same location does not move tail": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(0, 0),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     0,
				Visited: map[string]int{"0,0": 1},
			},
		},
		"Head and tail in same row and 1 column apart does not move tail": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(0, 1),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     1,
				Visited: map[string]int{"0,1": 1},
			},
		},
		"Head and tail in same column and 1 row apart does not move tail": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 1),
				Tail: rope.NewRopeEnd(0, 2),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     2,
				Visited: map[string]int{"0,2": 1},
			},
		},
		"Head and tail touching on single corner does not move tail": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(1, 1),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     1,
				Visited: map[string]int{"1,1": 1},
			},
		},
		"Head west of tail on same row and more than 1 column apart moves tail towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, -2),
				Tail: rope.NewRopeEnd(0, 0),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     -1,
				Visited: map[string]int{"0,0": 1, "0,-1": 1},
			},
		},
		"Head east of tail on same row and more than 1 column apart moves tail towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 2),
				Tail: rope.NewRopeEnd(0, 0),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     1,
				Visited: map[string]int{"0,0": 1, "0,1": 1},
			},
		},
		"Head north of tail on same column and more than 1 row apart moves tail towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(-2, 0),
			},
			want: &rope.RopeEnd{
				Row:     -1,
				Col:     0,
				Visited: map[string]int{"-2,0": 1, "-1,0": 1},
			},
		},
		"Head south of tail on same column and more than 1 row apart moves tail towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(1, 1),
				Tail: rope.NewRopeEnd(3, 1),
			},
			want: &rope.RopeEnd{
				Row:     2,
				Col:     1,
				Visited: map[string]int{"3,1": 1, "2,1": 1},
			},
		},
		"Head northwest of tail and more than 1 column apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(1, 0),
				Tail: rope.NewRopeEnd(0, 2),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     1,
				Visited: map[string]int{"0,2": 1, "1,1": 1},
			},
		},
		"Head northeast of tail and more than 1 column apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(1, 2),
				Tail: rope.NewRopeEnd(0, 0),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     1,
				Visited: map[string]int{"0,0": 1, "1,1": 1},
			},
		},
		"Head southeast of tail and more than 1 column apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 2),
				Tail: rope.NewRopeEnd(1, 0),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     1,
				Visited: map[string]int{"0,1": 1, "1,0": 1},
			},
		},
		"Head southwest of tail and more than 1 column apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(1, 2),
			},
			want: &rope.RopeEnd{
				Row:     0,
				Col:     1,
				Visited: map[string]int{"1,2": 1, "0,1": 1},
			},
		},
		"Head northwest of tail and more than 1 row apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(2, 0),
				Tail: rope.NewRopeEnd(0, 1),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     0,
				Visited: map[string]int{"0,1": 1, "1,0": 1},
			},
		},
		"Head northeast of tail and more than 1 row apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(2, 2),
				Tail: rope.NewRopeEnd(0, 1),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     2,
				Visited: map[string]int{"0,1": 1, "1,2": 1},
			},
		},
		"Head southwest of tail and more than 1 row apart moves tail diagonally towards head": {
			input: &rope.Rope{
				Head: rope.NewRopeEnd(0, 0),
				Tail: rope.NewRopeEnd(2, 1),
			},
			want: &rope.RopeEnd{
				Row:     1,
				Col:     0,
				Visited: map[string]int{"2,1": 1, "1,0": 1},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.input.UpdateTail()
			got := tc.input.Tail
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestRopeEnd_Move(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		numRows int
		numCols int
		want    *rope.RopeEnd
	}{
		"Moving by a positive number of rows updates rope end as expected": {
			numRows: 1,
			numCols: 0,
			want: &rope.RopeEnd{
				Row:     1,
				Col:     0,
				Visited: map[string]int{"0,0": 1, "1,0": 1},
			},
		},
		"Moving by a negative number of rows updates rope end as expected": {
			numRows: -1,
			numCols: 0,
			want: &rope.RopeEnd{
				Row:     -1,
				Col:     0,
				Visited: map[string]int{"0,0": 1, "-1,0": 1},
			},
		},
		"Moving by a positive number of columns updates rope end as expected": {
			numRows: 0,
			numCols: 1,
			want: &rope.RopeEnd{
				Row:     0,
				Col:     1,
				Visited: map[string]int{"0,0": 1, "0,1": 1},
			},
		},
		"Moving by a negative number of columns updates rope end as expected": {
			numRows: 0,
			numCols: -1,
			want: &rope.RopeEnd{
				Row:     0,
				Col:     -1,
				Visited: map[string]int{"0,0": 1, "0,-1": 1},
			},
		},
		"Moving to a position that has already been visited updates the map of visited positions": {
			numRows: 0,
			numCols: 0,
			want: &rope.RopeEnd{
				Row:     0,
				Col:     0,
				Visited: map[string]int{"0,0": 2},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := rope.NewRopeEnd(0, 0)
			r.Move(tc.numRows, tc.numCols)
			if !cmp.Equal(tc.want, r) {
				t.Error(cmp.Diff(tc.want, r))
			}
		})
	}
}

func TestHeadMovementFromLineErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Line with no quantity field returns error":            "U",
		"Line with invalid movement direction returns error":   "T 10",
		"Line with non-numerical quantity field returns error": "U ten",
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, _, err := rope.HeadMovementFromLine(tc)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestHeadMovementFromLine(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input       string
		wantNumRows int
		wantNumCols int
	}{
		"Valid up movement line returns expected number of rows": {
			input:       "U 10",
			wantNumRows: 10,
			wantNumCols: 0,
		},
		"Valid down movement line returns expected number of rows": {
			input:       "D 10",
			wantNumRows: -10,
			wantNumCols: 0,
		},
		"Valid left movement line returns expected number of columns": {
			input:       "L 10",
			wantNumRows: 0,
			wantNumCols: -10,
		},
		"Valid right movement line returns expected number of columns": {
			input:       "R 10",
			wantNumRows: 0,
			wantNumCols: 10,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotNumRows, gotNumCols, err := rope.HeadMovementFromLine(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantNumRows != gotNumRows {
				t.Errorf("want %d rows, got %d rows", tc.wantNumRows, gotNumRows)
			}
			if tc.wantNumCols != gotNumCols {
				t.Errorf("want %d cols, got %d cols", tc.wantNumCols, gotNumCols)
			}
		})
	}
}
