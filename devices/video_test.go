package devices_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/devices"
	"github.com/google/go-cmp/cmp"
)

func TestCrtScreen_WritePixel(t *testing.T) {
	t.Parallel()
	const (
		numRows = 2
		numCols = 4
	)
	testCases := map[string]struct {
		inputCpuCycle int
		inputPixelVal string
		want          *devices.CrtScreen
	}{
		"Writing a pixel at CPU cycle smaller than number of matrix columns writes pixel at expected location": {
			inputCpuCycle: 1,
			inputPixelVal: "#",
			want: &devices.CrtScreen{
				[]string{"#", "", "", ""},
				[]string{"", "", "", ""},
			},
		},
		"Writing a pixel at CPU cycle equal to number of matrix columns writes pixel at expected location": {
			inputCpuCycle: numCols,
			inputPixelVal: "#",
			want: &devices.CrtScreen{
				[]string{"", "", "", "#"},
				[]string{"", "", "", ""},
			},
		},
		"Writing a pixel at CPU cycle greater than number of matrix columns writes pixel at expected location": {
			inputCpuCycle: 7,
			inputPixelVal: "#",
			want: &devices.CrtScreen{
				[]string{"", "", "", ""},
				[]string{"", "", "#", ""},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := devices.NewCrtScreen(numRows, numCols)
			if err != nil {
				t.Fatal(err)
			}
			err = got.WritePixel(tc.inputCpuCycle, tc.inputPixelVal)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestCrtScreen_WritePixelErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		inputScreen   *devices.CrtScreen
		inputCpuCycle int
	}{
		"writing to a nil CRT screen returns error": {
			inputScreen: nil,
		},
		"writing to an out of bounds location on the CRT screen returns error": {
			inputScreen: &devices.CrtScreen{
				[]string{"", ""},
				[]string{"", ""},
			},
			inputCpuCycle: 100,
		},
		"writing to an empty CRT screen returns error": {
			inputScreen: &devices.CrtScreen{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.inputScreen.WritePixel(tc.inputCpuCycle, "#")
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestCrtScreen_OutputGivenAValidCrtScreenReturnsExpectedOutput(t *testing.T) {
	t.Parallel()
	screen := &devices.CrtScreen{
		[]string{".", ".", ".", "#"},
		[]string{"#", "#", ".", "#"},
	}
	want := `...#
##.#`
	got := screen.Output()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestDrawOnScreenWithValidInstructionsReturnsExpectedScreenOutput(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/valid-video-cpu-instructions.txt")
	if err != nil {
		t.Fatal(err)

	}
	defer f.Close()
	want := `##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....`
	got, err := devices.DrawOnScreen(f)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

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
