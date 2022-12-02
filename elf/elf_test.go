package elf_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/elf"
	"github.com/google/go-cmp/cmp"
)

func TestMaxCalories(t *testing.T) {
	t.Parallel()
	data := strings.NewReader(`1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
`)
	want := 24000
	got, err := elf.MaxCalories(data)

	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestTopCaloryCarriers(t *testing.T) {
	t.Parallel()
	data := strings.NewReader(`1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
`)
	want := []int{24000, 11000, 10000}
	got, err := elf.TopCaloryCarriers(data)

	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestCaloryStats_Insert(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		stats elf.CaloryStats
		input int
		want  elf.CaloryStats
	}{
		"Highest value is inserted at first position in highest counts slice": {
			stats: elf.CaloryStats{HighestCounts: []int{20, 15, 10}},
			input: 25,
			want:  elf.CaloryStats{HighestCounts: []int{25, 20, 15}},
		},
		"Mid-highest value is inserted at middle position in highest counts slice": {
			stats: elf.CaloryStats{HighestCounts: []int{20, 15, 10}},
			input: 18,
			want:  elf.CaloryStats{HighestCounts: []int{20, 18, 15}},
		},
		"Smallest value is inserted at last position in highest counts slice": {
			stats: elf.CaloryStats{HighestCounts: []int{20, 15, 10}},
			input: 12,
			want:  elf.CaloryStats{HighestCounts: []int{20, 15, 12}},
		},
		"Non-maximum value is not inserted into the highest counts slice": {
			stats: elf.CaloryStats{HighestCounts: []int{20, 15, 10}},
			input: 5,
			want:  elf.CaloryStats{HighestCounts: []int{20, 15, 10}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stats.Insert(tc.input)
			if !cmp.Equal(tc.want, tc.stats) {
				t.Error(cmp.Diff(tc.want, tc.stats))
			}
		})
	}
}
