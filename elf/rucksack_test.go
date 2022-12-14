package elf_test

import (
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/elf"
	"github.com/google/go-cmp/cmp"
)

func TestFindDuplicateRucksackItems(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  []rune
	}{
		"Input with single duplicate returns single duplicate": {
			input: "vJrwpWtwJgWrhcsFMMfFFhFp",
			want:  []rune{'p'},
		},
		"Input of length 2 with duplicates returns expected duplicate": {
			input: "JJ",
			want:  []rune{'J'},
		},
		"Input of length 2 with no duplicates returns nil slice": {
			input: "AB",
			want:  nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := elf.FindDuplicateRucksackItems(tc.input)
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(string(tc.want), string(got)))
			}
		})
	}
}

func TestSumDuplicateRucksackItemPriorities(t *testing.T) {
	t.Parallel()
	input := strings.NewReader(`vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
`)
	want := 157
	got, err := elf.SumDuplicateRucksackItemPriorities(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestFindBadgeInGroup(t *testing.T) {
	t.Parallel()
	input := [][]rune{
		[]rune("wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn"),
		[]rune("ttgJtRGJQctTZtZT"),
		[]rune("CrZsJsPPZsGzwwsLwLmpwMDw"),
	}
	want := 'Z'
	got, err := elf.FindBadgeInGroup(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if want != got {
		t.Errorf("want %s, got %s", string(want), string(got))
	}
}

func TestSumBadgeItemPriorities(t *testing.T) {
	t.Parallel()
	input := strings.NewReader(`vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
`)
	want := 70
	got, err := elf.SumBadgeItemPriorities(input)
	if err != nil {
		t.Fatal("got unexpected error: ", err)
	}
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
