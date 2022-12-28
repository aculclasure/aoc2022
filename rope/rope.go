// Package rope contains types and functions for performing rope movement
// simulations.
package rope

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Rope represents a rope with a head and tail rope end.
type Rope struct {
	Head *RopeEnd
	Tail *RopeEnd
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

// UpdateTail evaluates the position of the rope's head and tail ends and moves
// the tail end if the head and tail ends are no longer connected on any side or
// corner.
func (r *Rope) UpdateTail() {
	rowDiff := r.Head.Row - r.Tail.Row
	colDiff := r.Head.Col - r.Tail.Col
	switch {
	case rowDiff == 0 && colDiff == 0:
	case rowDiff == 0 && math.Abs(float64(colDiff)) > 1:
		if colDiff < 0 {
			colDiff++
		}
		if colDiff > 0 {
			colDiff--
		}
		r.Tail.Move(0, colDiff)
	case colDiff == 0 && math.Abs(float64(rowDiff)) > 1:
		if rowDiff < 0 {
			rowDiff++
		}
		if rowDiff > 0 {
			rowDiff--
		}
		r.Tail.Move(rowDiff, 0)
	case math.Abs(float64(colDiff)) > 1:
		if colDiff < 0 {
			colDiff++
		}
		if colDiff > 0 {
			colDiff--
		}
		r.Tail.Move(rowDiff, colDiff)
	case math.Abs(float64(rowDiff)) > 1:
		if rowDiff < 0 {
			rowDiff++
		}
		if rowDiff > 0 {
			rowDiff--
		}
		r.Tail.Move(rowDiff, colDiff)
	}
}

// NewRope returns a Rope struct with its head and tail ends initialized in the
// the same coordinate (row 0, column 0).
func NewRope() *Rope {
	return &Rope{
		Head: &RopeEnd{Row: 0, Col: 0, Visited: map[string]int{"0,0": 1}},
		Tail: &RopeEnd{Row: 0, Col: 0, Visited: map[string]int{"0,0": 1}},
	}
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
		return
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
// and applies these movements to the head end of a Rope. The starting point of
// the head and tail ends of the rope is row 0, column 0. After the movements
// are applied to the head end of the rope, the Rope struct is returned. An
// error is returned if the instructions argument is nil, if a movement line is
// invalidly formatted, or if there is a problem reading from the instructions.
func Run(instructions io.Reader) (*Rope, error) {
	if instructions == nil {
		return nil, errors.New("instructions argument must be non-nil")
	}
	rp := NewRope()
	scn := bufio.NewScanner(instructions)
	for scn.Scan() {
		line := scn.Text()
		numRows, numCols, err := HeadMovementFromLine(line)
		if err != nil {
			return nil, err
		}
		rp.MoveHead(numRows, numCols)
	}
	err := scn.Err()
	if err != nil {
		return nil, err
	}
	return rp, nil
}
