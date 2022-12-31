package devices

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type cpuInstruction struct {
	completesOnCycleNum int
}

func (c cpuInstruction) isComplete(currentCycleNum int) bool {
	return currentCycleNum >= c.completesOnCycleNum
}

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
