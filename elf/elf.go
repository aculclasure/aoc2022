package elf

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// TopCaloryCounts accepts an io.Reader pointing to a data set of how many
// calories are carried by each elf and returns a slice of ints representing
// the top 3 calory counts in ascending order. An error is returned if there
// is a problem reading the data set.
func TopCaloryCounts(data io.Reader) ([]int, error) {
	stats := NewCaloryStats()
	err := stats.ReadData(data)
	if err != nil {
		return nil, err
	}

	return stats.HighestCounts[:], nil
}

// CaloryStats represents statistics about elfs carrying calories.
type CaloryStats struct {
	HighestCounts []int
}

// Insert accepts a calory count and inserts it into appropriate position in
// the HighestCounts slice in the receiver.
func (c *CaloryStats) Insert(caloryCount int) {
	for i, v := range c.HighestCounts {
		if caloryCount <= v {
			continue
		}
		for k := len(c.HighestCounts) - 1; k >= i+1; {
			c.HighestCounts[k] = c.HighestCounts[k-1]
			k -= 1
		}
		c.HighestCounts[i] = caloryCount
		break
	}
}

// ReadData accepts an io.Reader pointing to a data set of how many
// calories are carried by each elf, parses it, and fills the HighestCounts
// slice in the receiver with the parsed data. An error is returned if there is
// a problem reading from the data source.
func (c *CaloryStats) ReadData(data io.Reader) error {
	if data == nil {
		return errors.New("data must point to a non-nil data source")
	}
	rdr := bufio.NewReader(data)
	currentElfCalories := 0
	for {
		line, err := rdr.ReadString('\n')
		if err != nil && err == io.EOF {
			c.Insert(currentElfCalories)
			break
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			c.Insert(currentElfCalories)
			currentElfCalories = 0
			continue
		}
		numCalories, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		currentElfCalories += numCalories

	}

	return nil
}

// NewCaloryStats returns a CaloryStats struct with an initialized HighestCounts
// field.
func NewCaloryStats() *CaloryStats {
	return &CaloryStats{HighestCounts: make([]int, 3)}
}
