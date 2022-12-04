package camp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CleaningAssignment struct {
	StartSector int
	EndSector   int
}

type CleaningPair struct {
	First  CleaningAssignment
	Second CleaningAssignment
}

func FullOverlapExists(pair CleaningPair) bool {
	switch {
	case pair.First.StartSector >= pair.Second.StartSector && pair.First.EndSector <= pair.Second.EndSector:
		return true
	case pair.Second.StartSector >= pair.First.StartSector && pair.Second.EndSector <= pair.First.EndSector:
		return true
	default:
		return false
	}
}

func PairFromInputLine(input string) (CleaningPair, error) {
	assignmentFields := strings.Split(input, ",")
	if len(assignmentFields) != 2 {
		return CleaningPair{}, fmt.Errorf("input must contain 2 assignments separated by a comma (got %d assignments for input %s)", len(assignmentFields), input)
	}

	var assignments []CleaningAssignment
	for _, asg := range assignmentFields {
		sectorFields := strings.Split(asg, "-")
		if len(sectorFields) != 2 {
			return CleaningPair{}, fmt.Errorf("assignment must be in the form start-end (got %s)", asg)
		}
		start, err := strconv.Atoi(sectorFields[0])
		if err != nil {
			return CleaningPair{}, err
		}
		end, err := strconv.Atoi(sectorFields[1])
		if err != nil {
			return CleaningPair{}, err
		}
		assignments = append(assignments, CleaningAssignment{StartSector: start, EndSector: end})
	}

	return CleaningPair{First: assignments[0], Second: assignments[1]}, nil
}

func GetFullyOverlappingPairs(schedules io.Reader) ([]CleaningPair, error) {
	if schedules == nil {
		return nil, errors.New("schedules must be a non-nil argument")
	}

	var fullyOverlapping []CleaningPair
	sc := bufio.NewScanner(schedules)
	for sc.Scan() {
		pair, err := PairFromInputLine(sc.Text())
		if err != nil {
			return nil, err
		}

		if FullOverlapExists(pair) {
			fullyOverlapping = append(fullyOverlapping, pair)
		}
	}
	err := sc.Err()
	if err != nil {
		return nil, err
	}

	return fullyOverlapping, nil
}
