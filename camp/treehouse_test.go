package camp_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/camp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	validTreeHeights   = []string{"30373", "25512", "65332", "33549", "35390"}
	invalidTreeHeights = []string{"3037a", "2B512", "653CC", "DDDDD", "12efg"}
)

func TestTreesFromBytesWithValidInputReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	input := []byte(`4043101133
1200221310
1144440434
2122404145
3241441301
`)
	want := []string{"4043101133", "1200221310", "1144440434", "2122404145", "3241441301"}
	got := camp.TreesFromBytes(input)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestTreesFromBytesWithNilInputReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	var input []byte
	want := []string{}
	got := camp.TreesFromBytes(input)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMaxScenicScoreWithValidTreeGridReturnsExpectedScore(t *testing.T) {
	t.Parallel()
	want := 8
	got, err := camp.MaxScenicScore(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestMaxScenicScoreWithInvalidTreeGridReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.MaxScenicScore(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}
func TestScenicScore(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		coordinate string
		want       int
	}{
		"Valid non-edge coordinate returns expected score": {
			coordinate: "3 2",
			want:       8,
		},
		"A different valid non-edge coordinate returns expected score": {
			coordinate: "1 2",
			want:       4,
		},
		"Coordinate on column grid edge returns 0 score": {
			coordinate: "1 0",
			want:       0,
		},
		"Coordinate on row grid edge returns 0 score": {
			coordinate: "0 1",
			want:       0,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := camp.ScenicScore(validTreeHeights, tc.coordinate)
			if err != nil {
				t.Fatal(err)
			}
			if tc.want != got {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}
func TestScenicScoreErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		trees      []string
		coordinate string
	}{
		"Coordinate with invalid row number returns error": {
			trees:      validTreeHeights,
			coordinate: "-1 2",
		},
		"Coordinate with invalid column number returns error": {
			trees:      validTreeHeights,
			coordinate: "0 100",
		},
		"Coordinate with non-numerical coordinate values returns error": {
			trees:      validTreeHeights,
			coordinate: "a b",
		},
		"Trees input with non-numerical height value returns error": {
			trees:      invalidTreeHeights,
			coordinate: "1 1",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := camp.ScenicScore(tc.trees, tc.coordinate)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestAllVisibleTreesWithValidInputReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	want := []string{
		"0 0", "0 1", "0 2", "0 3", "0 4",
		"1 0", "1 1", "1 2", "1 4",
		"2 0", "2 1", "2 3", "2 4",
		"3 0", "3 2", "3 4",
		"4 0", "4 1", "4 2", "4 3", "4 4"}
	got, err := camp.AllVisibleTrees(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	opt := cmpopts.SortSlices(func(i, j string) bool {
		return i < j
	})
	if !cmp.Equal(want, got, opt) {
		t.Error(cmp.Diff(want, got, opt))
	}
}

func TestAllVisibleTreesWithNonNumericalInputReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.AllVisibleTrees(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}

func TestVisibleFromLeftWithValidHeightsReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	want := []string{"0 0", "0 3", "1 0", "1 1", "2 0", "3 0", "3 2", "3 4", "4 0", "4 1", "4 3"}
	got, err := camp.VisibleFromLeft(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestVisibleFromLeftWithNonNumericalHeightsDataReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.VisibleFromLeft(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}

func TestVisibleFromRightWithValidHeightsReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	want := []string{"0 4", "0 3", "1 4", "1 2", "2 4", "2 3", "2 1", "2 0", "3 4", "4 4", "4 3"}
	got, err := camp.VisibleFromRight(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestVisibleFromRightWithNonNumericalHeightsDataReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.VisibleFromRight(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}

func TestVisibleFromTopWithValidHeightsReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	want := []string{"0 0", "2 0", "0 1", "1 1", "0 2", "1 2", "0 3", "4 3", "0 4", "3 4"}
	got, err := camp.VisibleFromTop(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestVisibleFromTopWithNonNumericalHeightsDataReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.VisibleFromTop(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}

func TestVisibleFromBottomWithValidHeightsReturnsExpectedCoordinateSlice(t *testing.T) {
	t.Parallel()
	want := []string{"4 0", "2 0", "4 1", "4 2", "3 2", "4 3", "4 4", "3 4"}
	got, err := camp.VisibleFromBottom(validTreeHeights)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestVisibleFromBottomWithNonNumericalHeightsDataReturnsError(t *testing.T) {
	t.Parallel()
	_, err := camp.VisibleFromBottom(invalidTreeHeights)
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}
