package devices

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CrtScreen represents the cathode ray tube on the elf communication device.
type CrtScreen [][]string

// WritePixel accepts a CPU cycle number and a pixel value as input, determines
// the appropriate location on the screen from the CPU cycle number and writes
// the pixel value at that location on the CRT screen. An error is returned if
// the receiver c is nil, if c is empty, or if the CPU cycle number converts to
// an invalid position on the CRT screen.
func (c *CrtScreen) WritePixel(cpuCycle int, pixelVal string) error {
	if c == nil {
		return errors.New("receiver must be non-nil")
	}
	if len(*c) == 0 {
		return errors.New("receiver must be a non empty crt matrix")
	}
	rowLength := len((*c)[0])
	row, col := (cpuCycle-1)/rowLength, (cpuCycle-1)%rowLength
	if row > len(*c) {
		return fmt.Errorf("cannot insert pixel at row %d, col %d into matrix of size %dx%d", row, col, len((*c)), len((*c)[0]))
	}
	if col > (len((*c)[0])) {
		return fmt.Errorf("cannot insert pixel at row %d, col %d into matrix of size %dx%d", row, col, len((*c)), len((*c)[0]))
	}
	(*c)[row][col] = pixelVal
	return nil
}

// Output returns the output of the CRT screen as a string.
func (c CrtScreen) Output() string {
	var rows []string
	for _, row := range c {
		nextRow := strings.Join(row, "")
		rows = append(rows, nextRow)
	}
	return strings.Join(rows, "\n")
}

// NewCrtScreen accepts a number of rows and columns and returns a CrtScreen
// configured with the given dimensions. An error is returned if an invalid
// value is given for the number of rows or number of columns.
func NewCrtScreen(numRows, numCols int) (*CrtScreen, error) {
	if numRows < 1 || numCols < 1 {
		return nil, errors.New("must specify a positive value for number of rows and columns")
	}
	var screen CrtScreen
	for i := 0; i < numRows; i++ {
		var nextRow []string
		for j := 0; j < numCols; j++ {
			nextRow = append(nextRow, "")
		}
		screen = append(screen, nextRow)
	}
	return &screen, nil
}

// cpuInstruction represents an instruction that is given to the video CPU on the
// elf communication device. It holds a field that indicates at what CPU cycle
// the instruction will be complete.
type cpuInstruction struct {
	completesOnCycleNum int
}

// isComplete accepts a CPU cycle as an int and determines if the CPU instruction
// is complete at that cycle.
func (c cpuInstruction) isComplete(currentCycleNum int) bool {
	return currentCycleNum >= c.completesOnCycleNum
}

// DrawOnScreen accepts line-separated instructions to the device's video CPU,
// draws pixels on a CRT screen according to the instructions and returns the
// output that is seen on the CRT screen. An error is returned if the instructions
// argument is nil, if an invalid instruction line is encountered, or if there is
// a problem reading the instructions.
func DrawOnScreen(instructions io.Reader) (string, error) {
	if instructions == nil {
		return "", errors.New("instructions must be non-nil")
	}
	const (
		numCols = 40
		numRows = 6
	)
	currentCycleNum := 1
	regValue := 1
	screen, err := NewCrtScreen(numRows, numCols)
	if err != nil {
		return "", err
	}
	scn := bufio.NewScanner(instructions)
	for scn.Scan() {
		line := scn.Text()
		fields := strings.Fields(line)
		switch {
		case len(fields) < 1:
			return "", errors.New("instruction line must be non-empty")
		case fields[0] != "noop" && fields[0] != "addx":
			return "", fmt.Errorf("instruction line must start with noop or addx (got %s)", line)
		case fields[0] == "addx" && len(fields) < 2:
			return "", fmt.Errorf("addx instruction line must have an adjustment value (got %s)", line)
		case fields[0] == "noop":
			instr := cpuInstruction{completesOnCycleNum: currentCycleNum + 1}
			for !instr.isComplete(currentCycleNum) {
				pixelVal := "."
				currentPosition := (currentCycleNum - 1) % numCols
				if overlaps(regValue, currentPosition) {
					pixelVal = "#"
				}
				screen.WritePixel(currentCycleNum, pixelVal)
				currentCycleNum++
			}
		default:
			delta, err := strconv.Atoi(fields[1])
			if err != nil {
				return "", err
			}
			instr := cpuInstruction{completesOnCycleNum: currentCycleNum + 2}
			for !instr.isComplete(currentCycleNum) {
				pixelVal := "."
				currentPosition := (currentCycleNum - 1) % numCols
				if overlaps(regValue, currentPosition) {
					pixelVal = "#"
				}
				screen.WritePixel(currentCycleNum, pixelVal)
				currentCycleNum++
			}
			regValue += delta
		}
	}
	err = scn.Err()
	if err != nil {
		return "", err
	}
	return screen.Output(), nil
}

// SignalStrengths accepts line-separated instructions to the device's video CPU,
// computes the power cycles from this set of instructions and returns them as a
// slice of ints. An error is returned if the instructions argument is nil, if
// an invalid instruction line is encountered, or if there is a problem reading
// the instructions.
func SignalStrengths(instructions io.Reader) ([]int, error) {
	if instructions == nil {
		return nil, errors.New("instructions must be non-nil")
	}
	interval := 20
	currentCycleNum := 1
	regValue := 1
	var sigStrengths []int
	scn := bufio.NewScanner(instructions)
	for scn.Scan() {
		line := scn.Text()
		fields := strings.Fields(line)
		switch {
		case len(fields) < 1:
			return nil, errors.New("instruction line must be non-empty")
		case fields[0] != "noop" && fields[0] != "addx":
			return nil, fmt.Errorf("instruction line must start with noop or addx (got %s)", line)
		case fields[0] == "addx" && len(fields) < 2:
			return nil, fmt.Errorf("addx instruction line must have an adjustment value (got %s)", line)
		case fields[0] == "noop":
			instr := cpuInstruction{completesOnCycleNum: currentCycleNum + 1}
			for !instr.isComplete(currentCycleNum) {
				if currentCycleNum%interval == 0 {
					sigStrengths = append(sigStrengths, interval*regValue)
					interval += 40
				}
				currentCycleNum++
			}
		default:
			delta, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			instr := cpuInstruction{completesOnCycleNum: currentCycleNum + 2}
			for !instr.isComplete(currentCycleNum) {
				if currentCycleNum%interval == 0 {
					sigStrengths = append(sigStrengths, currentCycleNum*regValue)
					interval += 40
				}
				currentCycleNum++
			}
			regValue += delta
		}
	}
	err := scn.Err()
	if err != nil {
		return nil, err
	}
	return sigStrengths, nil
}

// overlaps accepts the current position on a CRT line that a pixel will be written
// at and the position of the middle pixel of a sprite and returns true if there
// is an overlap.
func overlaps(currentPosition, spritePosition int) bool {
	return currentPosition >= spritePosition-1 && currentPosition <= spritePosition+1
}
