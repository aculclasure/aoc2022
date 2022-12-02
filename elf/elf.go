package elf

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func MaxCalories(data io.Reader) (int, error) {
	rdr := bufio.NewReader(data)
	max := 0
	currentElfCalories := 0
	for {
		line, err := rdr.ReadString('\n')
		if err != nil && err == io.EOF {
			if currentElfCalories > max {
				max = currentElfCalories
			}
			break
		}
		if err != nil {
			return 0, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if currentElfCalories > max {
				max = currentElfCalories
			}
			currentElfCalories = 0
			continue
		}
		numCalories, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		currentElfCalories += numCalories

	}
	return max, nil
}

func TopCaloryCarriers(data io.Reader) ([]int, error) {
	rdr := bufio.NewReader(data)
	stats := CaloryStats{HighestCounts: make([]int, 3)}
	currentElfCalories := 0
	for {
		line, err := rdr.ReadString('\n')
		if err != nil && err == io.EOF {
			stats.Insert(currentElfCalories)
			break
		}
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			stats.Insert(currentElfCalories)
			currentElfCalories = 0
			continue
		}
		numCalories, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		currentElfCalories += numCalories

	}
	return stats.HighestCounts, nil
}

type CaloryStats struct {
	HighestCounts []int
}

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
