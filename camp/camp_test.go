package camp_test

import (
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/camp"
	"github.com/google/go-cmp/cmp"
)

func TestFullOverlapExists(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input camp.CleaningPair
		want  bool
	}{
		"Non-overlapping pair with smaller first assignment returns false": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   2,
				},
				Second: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   4,
				},
			},
			want: false,
		},
		"Non-overlapping pair with smaller second assignment returns false": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   4,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   2,
				},
			},
			want: false,
		},
		"Partially overlapping pair returns false": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
				Second: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   6,
				},
			},
			want: false,
		},
		"Identical pairs returns true": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
			},
			want: true,
		},
		"Fully overlapping pair returns true": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 2,
					EndSector:   4,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
			},
			want: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := camp.FullOverlapExists(tc.input)
			if tc.want != got {
				t.Errorf("want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestOverlapExists(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input camp.CleaningPair
		want  bool
	}{
		"Non-overlapping pair with smaller first assignment returns false": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   2,
				},
				Second: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   4,
				},
			},
			want: false,
		},
		"Non-overlapping pair with smaller second assignment returns false": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   4,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   2,
				},
			},
			want: false,
		},
		"Partially overlapping pair returns true": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
				Second: camp.CleaningAssignment{
					StartSector: 3,
					EndSector:   6,
				},
			},
			want: true,
		},
		"Identical pairs returns true": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
			},
			want: true,
		},
		"Fully overlapping pair returns true": {
			input: camp.CleaningPair{
				First: camp.CleaningAssignment{
					StartSector: 2,
					EndSector:   4,
				},
				Second: camp.CleaningAssignment{
					StartSector: 1,
					EndSector:   5,
				},
			},
			want: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := camp.OverlapExists(tc.input)
			if tc.want != got {
				t.Errorf("want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestPairFromInputLine_ErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Empty input line returns an error":                         "",
		"Input line missing comma separator returns an error":       "1-23-4",
		"Input line with all non-number fields returns error":       "a-b,c-d",
		"Input line with a single invalid assignment returns error": "1-2,4-",
		"Input line with non-ranged assignments returns error":      "12,34",
	}
	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := camp.PairFromInputLine(input)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestPairFromInputLineReturnsPairForValidInput(t *testing.T) {
	t.Parallel()
	input := "1-2,3-4"
	want := camp.CleaningPair{
		First: camp.CleaningAssignment{
			StartSector: 1,
			EndSector:   2,
		},
		Second: camp.CleaningAssignment{
			StartSector: 3,
			EndSector:   4,
		},
	}

	got, err := camp.PairFromInputLine(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetFullyOverlappingPairs(t *testing.T) {
	t.Parallel()
	input := strings.NewReader(`1-2,3-4
2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8
`)
	want := []camp.CleaningPair{
		{
			First:  camp.CleaningAssignment{StartSector: 2, EndSector: 8},
			Second: camp.CleaningAssignment{StartSector: 3, EndSector: 7},
		},
		{
			First:  camp.CleaningAssignment{StartSector: 6, EndSector: 6},
			Second: camp.CleaningAssignment{StartSector: 4, EndSector: 6},
		},
	}
	got, err := camp.GetFullyOverlappingPairs(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetOverlappingPairs(t *testing.T) {
	t.Parallel()
	input := strings.NewReader(`2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8
`)
	want := []camp.CleaningPair{
		{
			First:  camp.CleaningAssignment{StartSector: 5, EndSector: 7},
			Second: camp.CleaningAssignment{StartSector: 7, EndSector: 9},
		},
		{
			First:  camp.CleaningAssignment{StartSector: 2, EndSector: 8},
			Second: camp.CleaningAssignment{StartSector: 3, EndSector: 7},
		},
		{
			First:  camp.CleaningAssignment{StartSector: 6, EndSector: 6},
			Second: camp.CleaningAssignment{StartSector: 4, EndSector: 6},
		},
		{
			First:  camp.CleaningAssignment{StartSector: 2, EndSector: 6},
			Second: camp.CleaningAssignment{StartSector: 4, EndSector: 8},
		},
	}
	got, err := camp.GetOverlappingPairs(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
