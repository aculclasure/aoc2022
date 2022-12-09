package cargo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var rgx = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)

// Crate represents a crate on the elf cargo ship. Each crate contains an item
// and has a stack assignment.
type Crate struct {
	Item  rune
	Stack int
}

// Movement represents a command to move crates from one stack into a
// different stack.
type Movement struct {
	Quantity  int
	SrcStack  int
	DestStack int
}

// Layout represents a cargo layout on the elf cargo ship.
type Layout struct {
	Stacks []*Stack
}

// NewLayout returns an initialized Layout.
func NewLayout(numStacks int) (*Layout, error) {
	if numStacks < 1 {
		return nil, fmt.Errorf("number of stacks must be at least 1 (got %d)", numStacks)
	}

	stacks := make([]*Stack, numStacks+1)
	for i := 1; i <= numStacks; i++ {
		stacks[i] = NewStack()
	}
	return &Layout{Stacks: stacks}, nil
}

// AddCrate accepts a Crate struct and adds its item to the appropriate stack
// in the layout. An error is returned if the crate contains an invalid stack
// index.
func (l *Layout) AddCrate(c Crate) error {
	if c.Stack < 1 || c.Stack >= len(l.Stacks) {
		return fmt.Errorf("stack index in given crate must be between 1 and %d inclusive (got %d)", len(l.Stacks)-1, c.Stack)
	}

	l.Stacks[c.Stack].Push(c.Item)
	return nil
}

// Move accepts a Movement instruction and applies it to the stacks in the cargo
// layout. An error is returned if there is a problem executing the movement
// instruction (like trying to move from or to a non-existent stack in the
// layout, trying to move a quantity larger than the number of crates in the
// source stack, or trying to move a negative quantity).
func (l *Layout) Move(mv Movement) error {
	switch {
	case mv.SrcStack < 1 || mv.SrcStack >= len(l.Stacks):
		return fmt.Errorf("the source stack in the movement must be between 1 and %d inclusive (got %d)", len(l.Stacks)-1, mv.SrcStack)
	case mv.DestStack < 1 || mv.DestStack >= len(l.Stacks):
		return fmt.Errorf("the destination stack in the movement must be between 1 and %d inclusive (got %d)", len(l.Stacks)-1, mv.DestStack)
	case mv.Quantity < 0:
		return fmt.Errorf("quantity to move must be 0 or greater (got %d)", mv.Quantity)
	case mv.Quantity > l.Stacks[mv.SrcStack].Size():
		return fmt.Errorf("quantity to move (%d) must not be greater than size of source stack (%d)", mv.Quantity, l.Stacks[mv.SrcStack].Size())
	}

	for i := 0; i < mv.Quantity; i++ {
		item, ok := l.Stacks[mv.SrcStack].Pop()
		if !ok {
			return fmt.Errorf("src stack %d must contain at least %d items (got %d)", mv.SrcStack, mv.Quantity, i)
		}
		l.Stacks[mv.DestStack].Push(item)
	}

	return nil
}

func (l *Layout) MoveWithCrateMover9001(mv Movement) error {
	switch {
	case mv.SrcStack < 1 || mv.SrcStack >= len(l.Stacks):
		return fmt.Errorf("the source stack in the movement must be between 1 and %d inclusive (got %d)", len(l.Stacks)-1, mv.SrcStack)
	case mv.DestStack < 1 || mv.DestStack >= len(l.Stacks):
		return fmt.Errorf("the destination stack in the movement must be between 1 and %d inclusive (got %d)", len(l.Stacks)-1, mv.DestStack)
	case mv.Quantity < 0:
		return fmt.Errorf("quantity to move must be 0 or greater (got %d)", mv.Quantity)
	case mv.Quantity > l.Stacks[mv.SrcStack].Size():
		return fmt.Errorf("quantity to move (%d) must not be greater than size of source stack (%d)", mv.Quantity, l.Stacks[mv.SrcStack].Size())
	}
	poppedItems := make([]rune, mv.Quantity)
	for i := 0; i < mv.Quantity; i++ {
		item, ok := l.Stacks[mv.SrcStack].Pop()
		if !ok {
			return fmt.Errorf("src stack %d must contain at least %d items (got %d)", mv.SrcStack, mv.Quantity, i)
		}
		poppedItems[i] = item
	}
	for i := range poppedItems {
		if mv.Quantity > 1 {
			l.Stacks[mv.DestStack].Push(poppedItems[len(poppedItems)-1-i])
			continue
		}
		l.Stacks[mv.DestStack].Push(poppedItems[i])
	}
	return nil
}

func (l *Layout) InitializeFromCrateRows(crates [][]Crate) error {
	for i := len(crates) - 1; i >= 0; i-- {
		row := crates[i]
		for _, c := range row {
			err := l.AddCrate(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTopItems adds the top item of each stack in the layout into a string and
// returns the string. An empty string is returned if the stacks in the layout
// are all empty.
func (l *Layout) GetTopItems() string {
	if len(l.Stacks) < 2 {
		return ""
	}

	topItems := ""
	for _, stk := range l.Stacks[1:] {
		top, ok := stk.Peek()
		if ok {
			topItems += string(top)
		}
	}

	return topItems
}

// GetCrates accepts a line of crate item information in the form
// "[A] [B] [C] ..." where A, B, and C represent crate items and the columnar
// position of A, B, and C indicate what stack they belong to and returns a
// slice of Crate structs.
func GetCrates(line string) []Crate {
	const crateFieldLen = 4
	var (
		crates        []Crate
		numSpacesRead int
	)
	currentStack := 0
	for _, v := range line {
		switch {
		case unicode.IsSpace(v):
			numSpacesRead++
			continue
		case v == '[' || v == ']':
			continue
		default:
			numEmptyStacks := numSpacesRead / crateFieldLen
			currentStack += 1 + numEmptyStacks
			crates = append(crates, Crate{Item: v, Stack: currentStack})
			numSpacesRead = 0
		}
	}

	return crates
}

// MovementFromLine accepts a string in the form "move QUANTITY from SRCSTACK to DESTSTACK"
// where QUANTITY is a quantity of items, SRCSTACK is an integer representing
// the source stack, and DESTSTACK is an integer representing the destination
// stack and returns a Movement struct with the parsed data. An error is returned
// if the line cannot be properly parsed into a Movement struct.
func MovementFromLine(line string) (Movement, error) {
	submatches := rgx.FindAllStringSubmatch(line, -1)
	if submatches == nil {
		return Movement{}, fmt.Errorf("line must be in the form move <qty> from <srcstack> to <deststack> (got %s)", line)
	}

	qty, err := strconv.Atoi(submatches[0][1])
	if err != nil {
		return Movement{}, fmt.Errorf("quantity field in line must be a valid integer (got %s)", submatches[0][1])
	}
	src, err := strconv.Atoi(submatches[0][2])
	if err != nil {
		return Movement{}, fmt.Errorf("srcstack field in line must be a valid integer (got %s)", submatches[0][2])
	}
	dest, err := strconv.Atoi(submatches[0][3])
	if err != nil {
		return Movement{}, fmt.Errorf("deststack field in line must be a valid integer (got %s)", submatches[0][3])
	}

	return Movement{Quantity: qty, SrcStack: src, DestStack: dest}, nil
}

// LayoutFromData accepts an io.Reader pointing to an initial cargo layout and
// a series of movements to apply cargo layout and returns the top items from
// from the stacks in the final cargo layout as a string. An error is returned
// if there is a problem reading the data or an invalid movement is applied
// to the layout.
func LayoutFromData(data io.Reader) (*Layout, error) {
	if data == nil {
		return nil, errors.New("data must be non-nil")
	}

	var (
		crateRows [][]Crate
		layout    *Layout
		err       error
	)
	scn := bufio.NewScanner(data)
	for scn.Scan() {
		line := scn.Text()
		switch {
		case strings.Contains(line, "["):
			crateRows = append(crateRows, GetCrates(line))
		case strings.HasPrefix(strings.TrimSpace(line), "1"):
			numStacks := len(strings.Fields(line))
			layout, err = NewLayout(numStacks)
			if err != nil {
				return nil, fmt.Errorf("got error creating layout: %s", err)
			}
			err = layout.InitializeFromCrateRows(crateRows)
			if err != nil {
				return nil, fmt.Errorf("got error initializing layout from crate rows: %s", err)
			}
		case strings.HasPrefix(line, "move"):
			mv, err := MovementFromLine(line)
			if err != nil {
				return nil, fmt.Errorf("got error creating movement from line %s: %s", line, err)
			}
			err = layout.MoveWithCrateMover9001(mv)
			if err != nil {
				return nil, fmt.Errorf("got error applying movement to layout: %s", err)
			}
		}
	}
	if err := scn.Err(); err != nil {
		return nil, err
	}

	return layout, nil
}
