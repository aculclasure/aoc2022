// Package rope contains types and functions for performing rope movement
// simulations.
package rope

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Opt represents a functional option that can be passed in during a call to the
// NewRope() function. It returns an error if the setting cannot be applied to
// the Rope struct.
type Opt func(*Rope) error

// WithNumKnots accepts an integer representing how many knots the rope should
// have a returns an Opt that configures a Rope to have this many knots. An error
// is returned if n is smaller than 2 since a rope must have at least a head and
// tail knot.
func WithNumKnots(n int) Opt {
	return func(r *Rope) error {
		if n < 2 {
			return errors.New("must specify at least 2 for number of knots to account for head and tail knots")
		}
		r.NumKnots = n
		return nil
	}
}

// Rope represents a rope with a head knot, tail knot, and an arbitrary number
// of knots between the head and tail knots.
type Rope struct {
	Head     *RopeEnd
	Tail     *RopeEnd
	NumKnots int
	Knots    []*RopeEnd
}

// MoveHead accepts a number of row and number of columns and moves the rope's
// head end accordingly. After each move of the head end, the tail end is also
// moved if needed.
func (r *Rope) MoveHead(numRows, numCols int) {
	rowDir := 1
	colDir := 1
	if numRows < 0 {
		numRows *= -1
		rowDir = -1
	}
	if numCols < 0 {
		numCols *= -1
		colDir = -1
	}
	for i := 0; i < numRows; i++ {
		r.Head.Move(rowDir*1, 0)
		r.UpdateTail()
	}
	for i := 0; i < numCols; i++ {
		r.Head.Move(0, colDir*1)
		r.UpdateTail()
	}
}

// UpdateTail compares the position of each knot in the Rope to the position of
// it's child knot (the next knot in the Knots slice contained in r) and adjusts
// the position of the child knot if it does not touch its parent knot on any side
// or corner.
func (r *Rope) UpdateTail() {
	for i := 0; i+1 < len(r.Knots); i++ {
		parent := r.Knots[i]
		child := r.Knots[i+1]
		rowDiff := parent.Row - child.Row
		colDiff := parent.Col - child.Col
		if abs(rowDiff) > 1 || abs(colDiff) > 1 {
			child.Move(delta(rowDiff), delta(colDiff))
		}
	}
}

// NewRope accepts an optional slice of Opts and returns a Rope struct with at least
// a head and tail knot configured at the default starting position (row 0,
// column 0). If the WithNumKnots() functional option is passed in as an argument
// then the rope creates that many knots at the default starting position. An
// error is returned if there is a problem applying a functional option to the
// Rope struct.
func NewRope(opts ...Opt) (*Rope, error) {
	head := &RopeEnd{Row: 0, Col: 0, Visited: map[string]int{"0,0": 1}}
	tail := &RopeEnd{Row: 0, Col: 0, Visited: map[string]int{"0,0": 1}}
	rp := &Rope{
		Head:     head,
		Tail:     tail,
		NumKnots: 2,
	}
	for _, o := range opts {
		err := o(rp)
		if err != nil {
			return nil, err
		}
	}
	knots := []*RopeEnd{head}
	for i := 2; i < rp.NumKnots; i++ {
		knots = append(knots, &RopeEnd{Visited: map[string]int{"0,0": 1}})
	}
	knots = append(knots, tail)
	rp.Knots = knots
	return rp, nil
}

// RopeFromRopeEnds accepts an optional slice of RopeEnd structs and returns a
// Rope struct configured with those knots. The first RopeEnd argument represents
// the head of the rope and the last RopeEnd argument represents the tail. If the
// slice of RopeEnd structs has a length smaller than 2, then an error is returned
// since the Rope must have at least a head and tail knot.
func RopeFromRopeEnds(ropeEnds ...*RopeEnd) (*Rope, error) {
	if len(ropeEnds) < 2 {
		return nil, errors.New("must give at least 2 rope ends to account for head and tail")
	}
	var knots []*RopeEnd
	knots = append(knots, ropeEnds...)
	return &Rope{
		Head:     knots[0],
		Tail:     knots[len(knots)-1],
		NumKnots: len(knots),
		Knots:    knots,
	}, nil
}

// RopeEnd represents the end of a rope. It contains fields to indicate the rope
// end's position by row and coordinate and also contains a history of all
// positions that it has visited.
type RopeEnd struct {
	Row     int
	Col     int
	Visited map[string]int
}

// Move accepts a number of rows and number of columns to move, moves the rope
// end to that position, and adds the new position to the map of all visited
// positions.
func (r *RopeEnd) Move(rowDelta, colDelta int) {
	r.Row += rowDelta
	r.Col += colDelta
	key := strconv.Itoa(r.Row) + "," + strconv.Itoa(r.Col)
	_, ok := r.Visited[key]
	if !ok {
		r.Visited[key] = 0
	}
	r.Visited[key] += 1
}

// NewRopeEnd accepts a row and column as integers and returns a RopeEnd struct
// that is initialized to that position.
func NewRopeEnd(row, col int) *RopeEnd {
	visited := map[string]int{
		strconv.Itoa(row) + "," + strconv.Itoa(col): 1,
	}
	return &RopeEnd{
		Row:     row,
		Col:     col,
		Visited: visited,
	}
}

// HeadMovementFromLine accepts a string representing a movement command:
// "U 4" (up 4)
// "D 4" (down 4)
// "L 4" (left 4)
// "R 4" (right 4)
// and parses this string into a number of rows and number of columns to move.
// The sign of the returned number of rows value or number of columns value indicates the
// direction to move, with negative values indicating a move down for rows and a
// move left for columns. An error is returned if the line does not contain at least
// 2 fields, if the first field is not a valid direction letter, or if the quantity
// field is not a valid integer.
func HeadMovementFromLine(line string) (numRows, numCols int, err error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		err = fmt.Errorf(`line must be in the form "<U|D|R|L> <qty>" (got line "%s")`, line)
		return
	}
	qty, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, err
	}
	direction := fields[0]
	switch {
	case direction == "U":
		numRows += qty
	case direction == "D":
		numRows -= qty
	case direction == "L":
		numCols -= qty
	case direction == "R":
		numCols += qty
	default:
		err = fmt.Errorf("direction must be one of U, D, R, L (got %s)", direction)
	}
	return
}

// Run accepts an io.Reader pointing to line-separated movement instructions
// and an integer representing how many knots should be in the rope and applies
// these movements to the head end of a rope configured with numKnots-many knots.
// The starting point of all knots in the rope is row 0, column 0. After the movements
// are applied to the head end of the rope, the Rope struct is returned. An
// error is returned if the instructions argument is nil, if a movement line is
// invalidly formatted, if there is a problem reading from the instructions, or
// if an invalid value is given for the numKnots argument (any value smaller than
// 2).
func Run(instructions io.Reader, numKnots int) (*Rope, error) {
	if instructions == nil {
		return nil, errors.New("instructions argument must be non-nil")
	}
	rp, err := NewRope(WithNumKnots(numKnots))
	if err != nil {
		return nil, err
	}
	scn := bufio.NewScanner(instructions)
	for scn.Scan() {
		line := scn.Text()
		numRows, numCols, err := HeadMovementFromLine(line)
		if err != nil {
			return nil, err
		}
		rp.MoveHead(numRows, numCols)
	}
	err = scn.Err()
	if err != nil {
		return nil, err
	}
	return rp, nil
}

// abs accepts an integer and returns its absolute value.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// delta accepts an integer and returns a delta value based on the sign of x.
func delta(x int) int {
	switch {
	case x == 0:
		return 0
	case x < 0:
		return -1
	default:
		return 1
	}
}
