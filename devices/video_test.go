package devices_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/devices"
	"github.com/google/go-cmp/cmp"
)

func TestSignalStrengthsWithValidInstructionsReturnsExpectedSignalStrengthSlice(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/valid-video-cpu-instructions.txt")
	if err != nil {
		t.Fatal(err)

	}
	defer f.Close()
	want := []int{420, 1140, 1800, 2940, 2880, 3960}
	got, err := devices.SignalStrengths(f)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSignalStrengthsErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]io.Reader{
		"nil instructions argument returns error": nil,
		"empty instruction line returns error":    strings.NewReader("\n"),
		"line that does not start with addx or noop returns error": strings.NewReader(
			"badinstruction -1\n",
		),
		"addx line missing adjustment value returns error":            strings.NewReader("addx\n"),
		"addx line with non-numerical adjustment value returns error": strings.NewReader("addx ten"),
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := devices.SignalStrengths(tc)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}
