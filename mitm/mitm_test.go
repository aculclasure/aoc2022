package mitm_test

import (
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/ds"
	"github.com/aculclasure/aoc2022/mitm"
	"github.com/google/go-cmp/cmp"
)

func TestMonkeysFromInputWithValidInputReturnsExpectedSliceOfMonkeys(t *testing.T) {
	t.Parallel()
	monkeyComp := cmp.Comparer(func(m1, m2 *mitm.Monkey) bool {
		if m1 == nil && m2 == nil {
			return true
		}
		if m1 == nil || m2 == nil {
			return false
		}
		if m1.ID != m2.ID {
			return false
		}
		if !(m1.Items == nil && m2.Items == nil) {
			if m1.Items == nil || m2.Items == nil {
				return false
			}
			if !cmp.Equal(m1.Items.PeekAllItems(), m2.Items.PeekAllItems()) {
				return false
			}
		}
		if !(m1.WorryCalc == nil && m2.WorryCalc == nil) {
			if m1.WorryCalc == nil || m2.WorryCalc == nil {
				return false
			}
			testValues := []int{-11, -5, 1, 5, 11}
			for _, v := range testValues {
				gotM1Res := m1.WorryCalc(v)
				gotM2Res := m2.WorryCalc(v)
				if gotM1Res != gotM2Res {
					return false
				}
			}
		}
		if m1.TestDivisor != m2.TestDivisor {
			return false
		}
		if m1.DestIfTrue != m2.DestIfTrue {
			return false
		}
		if m1.DestIfFalse != m2.DestIfFalse {
			return false
		}
		if m1.NumItemsInspected != m2.NumItemsInspected {
			return false
		}
		return true
	})
	monkeySliceComp := cmp.Comparer(func(s1, s2 []*mitm.Monkey) bool {
		if s1 == nil && s2 == nil {
			return true
		}
		if len(s1) != len(s2) {
			return false
		}
		for i := 0; i < len(s1); i++ {
			if !cmp.Equal(s1[i], s2[i], monkeyComp) {
				return false
			}
		}
		return true
	})
	testCases := map[string]struct {
		input string
		want  []*mitm.Monkey
	}{
		"Input containing single monkey returns slice with expected single monkey": {
			input: `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
			  If true: throw to monkey 2
			  If false: throw to monkey 3

			`,
			want: []*mitm.Monkey{
				{
					ID:    0,
					Items: ds.NewQueueFromItems(79, 98),
					WorryCalc: func(old int) int {
						return old * 19
					},
					TestDivisor:       23,
					DestIfTrue:        2,
					DestIfFalse:       3,
					NumItemsInspected: 0,
				},
			},
		},
		"Input containing multiple monkeys returns slice with multiple monkeys": {
			input: `Monkey 0:
			  Starting items: 79, 98
			  Operation: new = old * 19
			  Test: divisible by 23
			    If true: throw to monkey 2
			    If false: throw to monkey 3

			Monkey 1:
                          Starting items: 54, 65, 75, 74
                          Operation: new = old + 6
                          Test: divisible by 19
                            If true: throw to monkey 2
                            If false: throw to monkey 0

			`,
			want: []*mitm.Monkey{
				{
					ID:    0,
					Items: ds.NewQueueFromItems(79, 98),
					WorryCalc: func(old int) int {
						return old * 19
					},
					TestDivisor:       23,
					DestIfTrue:        2,
					DestIfFalse:       3,
					NumItemsInspected: 0,
				},
				{
					ID:    1,
					Items: ds.NewQueueFromItems(54, 65, 75, 74),
					WorryCalc: func(old int) int {
						return old + 6
					},
					TestDivisor:       19,
					DestIfTrue:        2,
					DestIfFalse:       0,
					NumItemsInspected: 0,
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := mitm.MonkeysFromInput(strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.want, got, monkeySliceComp) {
				t.Error(cmp.Diff(tc.want, got, monkeySliceComp))
			}
		})
	}
}

func TestMonkeysFromInputErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Monkey ID line with no monkey id value returns error": `Monkey
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Monkey ID line with no monkey id value but multiple fields returns error": `Monkey : : :
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line of starting items with no starting items returns error": `Monkey 0:
			Starting items:
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line of starting items with invalid starting item returns error": `Monkey 0:
			Starting items: 79, 98, a
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line with incomplete operation returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new =
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line with undefined operation returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old times 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line with invalid operand in operation returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * nineteen
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line with invalid test condition returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"Line with invalid test condition operand returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by twenty-three
				If true: throw to monkey 2
			  	If false: throw to monkey 3
			`,
		"True test result line with missing destination monkey returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey
			  	If false: throw to monkey 3
			`,
		"True test result line with non-integer destination monkey returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey two
			  	If false: throw to monkey 3
			`,
		"True test result line with negative destination monkey value returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey -5
			  	If false: throw to monkey 3
			`,
		"False test result line with missing destination monkey returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey
			`,
		"False test result line with non-integer destination monkey returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey three
			`,
		"False test result line with negative destination monkey value returns error": `Monkey 0:
			Starting items: 79, 98
			Operation: new = old * 19
			Test: divisible by 23
				If true: throw to monkey 2
			  	If false: throw to monkey -3
			`,
	}
	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := mitm.MonkeysFromInput(strings.NewReader(input))
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestRunWithValidMonkeysInputUpdatesMonkeyItemsAsExpected(t *testing.T) {
	t.Parallel()
	numRounds := 1
	testCases := map[string]struct {
		input    []*mitm.Monkey
		adjuster mitm.WorryLevelAdjuster
		want     [][]int
	}{
		"Worry level divisor of 3 returns expected items after 1 round": {
			input:    getTestMonkeys(),
			adjuster: mitm.AdjustWorryLevelPart1{Divisor: 3},
			want: [][]int{
				{20, 23, 27, 26},
				{2080, 25, 167, 207, 401, 1046},
				nil,
				nil,
			},
		},
		"No worry level reduction effect returns expected items after 1 round": {
			input:    getTestMonkeys(),
			adjuster: mitm.AdjustWorryLevelPart2{CommonMultiple: mitm.CommonMultiple(getTestMonkeys())},
			want: [][]int{
				{60, 71, 81, 80},
				{77, 1504, 1865, 6244, 3603, 9412},
				nil,
				nil,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := mitm.Run(tc.input, numRounds, tc.adjuster)
			if err != nil {
				t.Fatal(err)
			}
			var got [][]int
			for _, mk := range tc.input {
				got = append(got, mk.Items.PeekAllItems())
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestRunWithValidMonkeysInputSetsNumItemsInspectedFieldAsExpected(t *testing.T) {
	t.Parallel()
	numRounds := 1
	testCases := map[string]struct {
		input    []*mitm.Monkey
		adjuster mitm.WorryLevelAdjuster
		want     []int
	}{
		"Worry level divisor of 3 sets number of items inspected for each monkey to expected values after 1 round": {
			input:    getTestMonkeys(),
			adjuster: mitm.AdjustWorryLevelPart1{Divisor: 3},
			want:     []int{2, 4, 3, 5},
		},
		"No worry level reduction effect sets number of items inspected for each monkey to expected values after 1 round": {
			input:    getTestMonkeys(),
			adjuster: mitm.AdjustWorryLevelPart2{CommonMultiple: mitm.CommonMultiple(getTestMonkeys())},
			want:     []int{2, 4, 3, 6},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := mitm.Run(tc.input, numRounds, tc.adjuster)
			if err != nil {
				t.Fatal(err)
			}
			var got []int
			for _, mk := range tc.input {
				got = append(got, mk.NumItemsInspected)
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestMonkeyBusinessWithValidMonkeysInputReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input     []*mitm.Monkey
		numRounds int
		adjuster  mitm.WorryLevelAdjuster
		want      int
	}{
		"Running 20 rounds with worry level reduction divisor of 3 returns expected value": {
			input:     getTestMonkeys(),
			numRounds: 20,
			adjuster:  mitm.AdjustWorryLevelPart1{Divisor: 3},
			want:      10605,
		},
		"Running 10000 rounds with no worry level reduction returns expected value": {
			input:     getTestMonkeys(),
			numRounds: 10000,
			adjuster:  mitm.AdjustWorryLevelPart2{CommonMultiple: mitm.CommonMultiple(getTestMonkeys())},
			want:      2713310158,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := mitm.Run(tc.input, tc.numRounds, tc.adjuster)
			if err != nil {
				t.Fatal(err)
			}
			for _, mk := range tc.input {
				t.Logf("Monkey %d inspected %d items after %d rounds\n", mk.ID, mk.NumItemsInspected, tc.numRounds)
			}
			got := mitm.MonkeyBusiness(tc.input)
			if tc.want != got {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}

func TestCommonMultipleWithValidInputReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	monkeys := getTestMonkeys()
	want := 96577
	got := mitm.CommonMultiple(monkeys)
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func getTestMonkeys() []*mitm.Monkey {
	return []*mitm.Monkey{
		{
			ID:    0,
			Items: ds.NewQueueFromItems(79, 98),
			WorryCalc: func(old int) int {
				return old * 19
			},
			TestDivisor: 23,
			DestIfTrue:  2,
			DestIfFalse: 3,
		},
		{
			ID:    1,
			Items: ds.NewQueueFromItems(54, 65, 75, 74),
			WorryCalc: func(old int) int {
				return old + 6
			},
			TestDivisor: 19,
			DestIfTrue:  2,
			DestIfFalse: 0,
		},
		{
			ID:    2,
			Items: ds.NewQueueFromItems(79, 60, 97),
			WorryCalc: func(old int) int {
				return old * old
			},
			TestDivisor: 13,
			DestIfTrue:  1,
			DestIfFalse: 3,
		},
		{
			ID:    3,
			Items: ds.NewQueueFromItems(74),
			WorryCalc: func(old int) int {
				return old + 3
			},
			TestDivisor: 17,
			DestIfTrue:  0,
			DestIfFalse: 1,
		},
	}
}
