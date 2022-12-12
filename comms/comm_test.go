package comms_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/comms"
)

func TestHasUniqueCharsWithEmptyInputReturnsError(t *testing.T) {
	t.Parallel()
	input := []rune("")
	_, err := comms.HasUniqueChars(input)
	if err == nil {
		t.Fatal("expected an error but did not get one")
	}
}

func TestHasUniqueChars(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input []rune
		want  bool
	}{
		"Input with single character returns true": {
			input: []rune("j"),
			want:  true,
		},
		"Input with unique characters returns true": {
			input: []rune("jpqm"),
			want:  true,
		},
		"Input with non-unique characters returns false": {
			input: []rune("jpqj"),
			want:  false,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := comms.HasUniqueChars(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if tc.want != got {
				t.Errorf("want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestStartPacketMarker(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  int
	}{
		"Input of length 4 with all unique characters returns 4": {
			input: "abcd",
			want:  4,
		},
		"Input of length 4 with non-unique characters returns -1": {
			input: "abca",
			want:  -1,
		},
		"Input with marker at position 5 returns 5": {
			input: "bvwbjplbgvbhsrlpgdmjqwftvncz",
			want:  5,
		},
		"Input with marker at position 11 returns 11": {
			input: "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
			want:  11,
		},
		"Input that is too small returns -1": {
			input: "abc",
			want:  -1,
		},
		"Input with length greater than 4 and marker at end of input returns expected index": {
			input: "zcfzfwzzqfr",
			want:  11,
		},
		"Input with length greater than 4 and no unique characters returns -1": {
			input: "aabbccddeeffgghhiijjkkllmm",
			want:  -1,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := comms.StartPacketMarker(tc.input)
			if tc.want != got {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}
