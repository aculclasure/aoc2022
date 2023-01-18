package mitm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/aculclasure/aoc2022/ds"
)

var worryFuncs = map[string]func(int, int) int{
	"*": func(a, b int) int { return a * b },
	"+": func(a, b int) int { return a + b },
}

// Monkey represents a monkey participating in the man-in-the-middle game.
type Monkey struct {
	ID                int
	Items             *ds.Queue[int]
	WorryCalc         func(int) int
	TestDivisor       int
	DestIfTrue        int
	DestIfFalse       int
	NumItemsInspected int
}

// AdjustWorryLevelPart1 provides a type with the Adjust behavior needed for
// solving part one of the problem.
type AdjustWorryLevelPart1 struct {
	Divisor int
}

// Adjust accepts a worry level as an int and returns the result of dividing the
// worry level by 3.
func (a AdjustWorryLevelPart1) Adjust(oldWorryLevel int) int {
	return oldWorryLevel / a.Divisor
}

// AdjustWorryLevelPart2 provides a type with the Adjust behavior needed for
// solving part two of the problem where there is no post-inspection reduction
// of worry levels.
type AdjustWorryLevelPart2 struct {
	CommonMultiple int
}

// Adjust accepts a worry level as an int and returns the modulus of the worry
// level with the common multiple field in the AdjustWorryLevelPart2 receiver.
func (a AdjustWorryLevelPart2) Adjust(oldWorryLevel int) int {
	return oldWorryLevel % a.CommonMultiple
}

// MonkeysFromInput accepts an io.Reader pointing to line-separated Monkey
// attribute data and returns a slice of Monkey structs built from those
// attributes. An error is returned if the input cannot be processed or if there
// is a problem constructing a Monkey struct from the given attributes.
func MonkeysFromInput(input io.Reader) ([]*Monkey, error) {
	if input == nil {
		return nil, errors.New("input must be non-nil")
	}
	var (
		monkeys []*Monkey
		next    *Monkey
	)
	scn := bufio.NewScanner(input)
	for scn.Scan() {
		line := strings.TrimSpace(scn.Text())
		if !strings.HasPrefix(line, "Monkey") {
			continue
		}
		next, err := parseMonkey(line, scn)
		if err != nil {
			return nil, err
		}
		monkeys = append(monkeys, next)
	}
	if err := scn.Err(); err != nil {
		return nil, err
	}
	if next != nil {
		monkeys = append(monkeys, next)
	}
	return monkeys, nil
}

// WorryLevelAdjuster provides an interface for the post-inspection worry level
// adjustment behavior. For part one of the problem, the worry level is adjusted
// by dividing it by 3. For part two of the problem, the worry level is adjusted
// by using a modulus.
type WorryLevelAdjuster interface {
	Adjust(int) int
}

// Run accepts a slice of Monkey structs and an number of rounds to execute and
// runs the monkey-in-the-middle game for the specified number of rounds. An
// error is returned if invalid arguments are provided to the function or if a
// monkey attempts to throw one of it's items to an invalid destination monkey.
func Run(monkeys []*Monkey, numRounds int, adjuster WorryLevelAdjuster) error {
	if monkeys == nil {
		return errors.New("monkeys argument must be non-nil")
	}
	if numRounds < 1 {
		return fmt.Errorf("number of rounds must be at least 1, got %d", numRounds)
	}
	if adjuster == nil {
		return errors.New("must provide a non-nil worry level adjuster argument")
	}
	for i := 0; i < numRounds; i++ {
		for _, mk := range monkeys {
			for mk.Items.Size() > 0 {
				worryLevel, _ := mk.Items.Dequeue()
				mk.NumItemsInspected++
				worryLevel = adjuster.Adjust(mk.WorryCalc(worryLevel))
				testPassed := (worryLevel % mk.TestDivisor) == 0
				destMonkey := mk.DestIfTrue
				if !testPassed {
					destMonkey = mk.DestIfFalse
				}
				if destMonkey >= len(monkeys) {
					return fmt.Errorf("destination monkey must be between 0 and %d inclusive, got %d", len(monkeys)-1, destMonkey)
				}
				monkeys[destMonkey].Items.Enqueue(worryLevel)
			}
		}
	}
	return nil
}

// MonkeyBusiness accepts a slice of Monkey structs and returns the product of the
// number of inspected items of the 2 busiest monkeys in the slice. Business
// is defined as how many items that monkey has inspected. An empty monkey slice
// returns 0.
func MonkeyBusiness(monkeys []*Monkey) int {
	if len(monkeys) == 0 {
		return 0
	}
	if len(monkeys) == 1 {
		return monkeys[0].NumItemsInspected
	}
	var copy []*Monkey
	copy = append(copy, monkeys...)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].NumItemsInspected < copy[j].NumItemsInspected
	})
	return copy[len(copy)-1].NumItemsInspected * copy[len(copy)-2].NumItemsInspected
}

// CommonMultiple accepts a slice of Monkey structs and returns the common
// multiple of all the TestDivisor fields in each monkey. 0 is returned if a nil
// slice of Monkeys is given.
func CommonMultiple(monkeys []*Monkey) int {
	if monkeys == nil {
		return 0
	}
	multiple := 1
	for _, m := range monkeys {
		multiple *= m.TestDivisor
	}
	return multiple
}

// parseMonkey accepts a monkey indicator line (e.g. "Monkey 0") and a Scanner
// that reads line-separated monkey attributes for the indicated monkey and
// returns a Monkey struct. An error is returned if there is problem scanning
// the input or if an invalid attribute is encountered.
func parseMonkey(monkeyIDLine string, scn *bufio.Scanner) (*Monkey, error) {
	flds := strings.Fields(monkeyIDLine)
	if len(flds) < 2 {
		return nil, fmt.Errorf(`expected monkey id line to have an id value, got "%s"`, monkeyIDLine)
	}
	id, err := strconv.Atoi(strings.TrimSuffix(flds[1], ":"))
	if err != nil {
		return nil, err
	}
	monkey := &Monkey{ID: id}
	for scn.Scan() {
		line := strings.TrimSpace(scn.Text())
		switch {
		case strings.HasPrefix(line, "Starting items: "):
			monkeyItems := ds.NewQueue[int]()
			res := strings.Split(line, "Starting items: ")
			if len(res) < 2 {
				return nil, fmt.Errorf("line must contain worry levels, got %s", line)
			}
			items := strings.Split(res[1], ",")
			if len(items) < 1 {
				return nil, fmt.Errorf("line must contain worry levels, got %s", line)
			}
			for _, item := range items {
				val, err := strconv.Atoi(strings.TrimSpace(item))
				if err != nil {
					return nil, err
				}
				monkeyItems.Enqueue(val)
			}
			monkey.Items = monkeyItems
		case strings.HasPrefix(line, "Operation:"):
			flds = strings.Fields(line)
			if len(flds) < 6 {
				return nil, fmt.Errorf("line must contain a valid worry calculation, got %s", line)
			}
			op := flds[4]
			worryCalc, ok := worryFuncs[op]
			if !ok {
				return nil, fmt.Errorf("no worry calculation function defined for operator %s", op)
			}
			opnd := flds[5]
			if opnd == "old" {
				monkey.WorryCalc = func(old int) int {
					return worryCalc(old, old)
				}
				continue
			}
			opndVal, err := strconv.Atoi(opnd)
			if err != nil {
				return nil, err
			}
			monkey.WorryCalc = func(old int) int {
				return worryCalc(old, opndVal)
			}
		case strings.HasPrefix(line, "Test:"):
			flds = strings.Fields(line)
			if len(flds) < 4 {
				return nil, fmt.Errorf("line must contain a valid test condition, got %s", line)
			}
			div, err := strconv.Atoi(flds[3])
			if err != nil {
				return nil, err
			}
			if div == 0 {
				return nil, errors.New("test condition line must contain a non-zero divisor")
			}
			monkey.TestDivisor = div
		case strings.HasPrefix(line, "If true:") || strings.HasPrefix(line, "If false:"):
			flds = strings.Fields(line)
			if len(flds) < 6 {
				return nil, fmt.Errorf("line must contain a valid test result definition, got %s", line)
			}
			dest, err := strconv.Atoi(flds[5])
			if err != nil {
				return nil, err
			}
			if dest < 0 {
				return nil, fmt.Errorf("destination must be a non-negative int, got %d", dest)
			}
			if flds[1] == "true:" {
				monkey.DestIfTrue = dest
			} else {
				monkey.DestIfFalse = dest
			}
		case line == "":
			return monkey, nil
		default:
			return nil, fmt.Errorf("no processing logic exists for line %s", line)
		}
	}
	err = scn.Err()
	if err != nil {
		return nil, err
	}
	return monkey, nil
}
