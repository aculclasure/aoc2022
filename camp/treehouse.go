package camp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// TreesFromBytes accepts a slice of bytes that is intended to come from a file
// containing line-separated tree height data and returns a slice where each string
// is a row of tree height data.
func TreesFromBytes(input []byte) []string {
	return strings.Fields(strings.TrimSpace(string(input)))
}

// MaxScenicStore accepts a slice of strings representing a grid of tree height
// data and returns the highest scenic score value for a single tree within the
// grid. An error is returned if the grid contains any non-numerical data.
func MaxScenicScore(trees []string) (int, error) {
	type scoreResult struct {
		score int
		err   error
	}
	var wg sync.WaitGroup
	results := make(chan scoreResult)
	for i := range trees {
		wg.Add(1)
		go func(rowIdx int) {
			defer wg.Done()
			maxScore := -1
			for colIdx := range trees[rowIdx] {
				score, err := ScenicScore(trees, strconv.Itoa(rowIdx)+" "+strconv.Itoa(colIdx))
				if err != nil {
					results <- scoreResult{err: err}
					return
				}
				if score > maxScore {
					maxScore = score
				}
			}
			results <- scoreResult{score: maxScore}
		}(i)
	}
	go func() {
		defer close(results)
		wg.Wait()
	}()
	maxScore := -1
	for res := range results {
		if res.err != nil {
			return 0, res.err
		}
		if res.score > maxScore {
			maxScore = res.score
		}
	}
	return maxScore, nil
}

// ScenicScore accepts a slice of strings representing a grid of tree height
// data and a coordinate and returns the scenic score for the tree at that given
// coordinate. An error is returned if the grid contains invalid data (e.g. non
// integer characters) or if the coordinate is invalid.
func ScenicScore(trees []string, coord string) (int, error) {
	numRows := len(trees)
	if numRows == 0 {
		return 0, errors.New("trees must be non-empty slice")
	}
	numCols := len(trees[0])
	if numCols == 0 {
		return 0, errors.New("trees must contain at least 1 row and column of height data")
	}
	coordFields := strings.Fields(coord)
	if len(coordFields) != 2 {
		return 0, fmt.Errorf(`coordinate must be in the form "r c" (got %s)`, coord)
	}
	row, err := strconv.Atoi(coordFields[0])
	if err != nil {
		return 0, fmt.Errorf("got error %s, is the row in your coordinate a valid integer?", err)
	}
	col, err := strconv.Atoi(coordFields[1])
	if err != nil {
		return 0, fmt.Errorf("got error %s, is the column in your coordinate a valid integer?", err)
	}
	if row < 0 || row >= numRows {
		return 0, fmt.Errorf("row must be a value from 0-%d (got %d)", numRows-1, row)
	}
	if col < 0 || col >= numCols {
		return 0, fmt.Errorf("col must be a value from 0-%d (got %d)", numCols-1, col)
	}
	if row == 0 || row == numRows-1 || col == 0 || col == numCols-1 { // Row or column is on an edge
		return 0, nil
	}
	startHeight, err := strconv.Atoi(string(trees[row][col]))
	if err != nil {
		return 0, fmt.Errorf("got error %s, check that your trees input only contains integers", err)
	}
	score := 1
	viewingDistance := 0
	for c := col; c+1 < numCols; c++ {
		viewingDistance++
		currentHeight, err := strconv.Atoi(string(trees[row][c+1]))
		if err != nil {
			return 0, fmt.Errorf("got error %s, check that your trees input only contains integers", err)
		}
		if currentHeight >= startHeight {
			break
		}
	}
	score *= viewingDistance
	viewingDistance = 0
	for c := col; c-1 >= 0; c-- {
		viewingDistance++
		currentHeight, err := strconv.Atoi(string(trees[row][c-1]))
		if err != nil {
			return 0, fmt.Errorf("got error %s, check that your trees input only contains integers", err)
		}
		if currentHeight >= startHeight {
			break
		}
	}
	score *= viewingDistance
	viewingDistance = 0
	for r := row; r+1 < numRows; r++ {
		viewingDistance++
		currentHeight, err := strconv.Atoi(string(trees[r+1][col]))
		if err != nil {
			return 0, fmt.Errorf("got error %s, check that your trees input only contains integers", err)
		}
		if currentHeight >= startHeight {
			break
		}
	}
	score *= viewingDistance
	viewingDistance = 0
	for r := row; r-1 >= 0; r-- {
		viewingDistance++
		currentHeight, err := strconv.Atoi(string(trees[r-1][col]))
		if err != nil {
			return 0, fmt.Errorf("got error %s, check that your trees input only contains integers", err)
		}
		if currentHeight >= startHeight {
			break
		}
	}
	score *= viewingDistance
	return score, nil
}

// AllVisibleTrees accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from outside the grid. Each coordinate in the returned slice is in the form
// "R C", where R indicates the row and C indicates the column of a visible tree
// (e.g. "0 3", "1 4", etc.) An error is returned if the grid contains any
// non-numerical data.
func AllVisibleTrees(trees []string) ([]string, error) {
	type result struct {
		visTrees []string
		err      error
	}
	var wg sync.WaitGroup
	results := make(chan result)
	visibleSet := make(map[string]struct{})
	wg.Add(4)
	go func() {
		defer wg.Done()
		visFromLeft, err := VisibleFromLeft(trees)
		results <- result{visTrees: visFromLeft, err: err}
	}()
	go func() {
		defer wg.Done()
		visFromTop, err := VisibleFromTop(trees)
		results <- result{visTrees: visFromTop, err: err}
	}()
	go func() {
		defer wg.Done()
		visFromRight, err := VisibleFromRight(trees)
		results <- result{visTrees: visFromRight, err: err}
	}()
	go func() {
		defer wg.Done()
		visFromBtm, err := VisibleFromBottom(trees)
		results <- result{visTrees: visFromBtm, err: err}
	}()
	go func() {
		defer close(results)
		wg.Wait()
	}()
	for r := range results {
		if r.err != nil {
			return nil, r.err
		}
		for _, visTreeCoord := range r.visTrees {
			visibleSet[visTreeCoord] = struct{}{}
		}
	}
	var visible []string
	for k := range visibleSet {
		visible = append(visible, k)
	}
	return visible, nil
}

// VisibleFromLeft accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from outside the grid on the left side. Each coordinate in the returned slice
// is in the form "R C", where R indicates the row and C indicates the column of
// a visible tree (e.g. "0 3", "1 4", etc.) An error is returned if the grid
// contains any non-numerical data.
func VisibleFromLeft(trees []string) ([]string, error) {
	var visible []string
	for row := 0; row < len(trees); row++ {
		maxHeight := -1
		for col := 0; col < len(trees[row]); col++ {
			currentHeight, err := strconv.Atoi(string(trees[row][col]))
			if err != nil {
				return nil, fmt.Errorf("got error %s, does your trees argument only contain numerical data?", err)
			}
			if currentHeight > maxHeight {
				maxHeight = currentHeight
				visible = append(visible, strconv.Itoa(row)+" "+strconv.Itoa(col))
			}
		}
	}
	return visible, nil
}

// VisibleFromRight accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from outside the grid on the right side. Each coordinate in the returned slice
// is in the form "R C", where R indicates the row and C indicates the column of
// a visible tree (e.g. "0 3", "1 4", etc.) An error is returned if the grid
// contains any non-numerical data.
func VisibleFromRight(trees []string) ([]string, error) {
	var visible []string
	for row := 0; row < len(trees); row++ {
		maxHeight := -1
		for col := len(trees[row]) - 1; col >= 0; col-- {
			currentHeight, err := strconv.Atoi(string(trees[row][col]))
			if err != nil {
				return nil, fmt.Errorf("got error %s, does your trees argument only contain numerical data?", err)
			}
			if currentHeight > maxHeight {
				maxHeight = currentHeight
				visible = append(visible, strconv.Itoa(row)+" "+strconv.Itoa(col))
			}
		}
	}
	return visible, nil
}

// VisibleFromTop accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from the outside on the top of the grid. Each coordinate in the returned
// slice is in the form "R C", where R indicates the row and C indicates the
// column of a visible tree (e.g. "0 3", "1 4", etc.) An error is returned if
// the grid contains any non-numerical data.
func VisibleFromTop(trees []string) ([]string, error) {
	if len(trees) == 0 {
		return nil, nil
	}
	var visible []string
	numCols := len(trees[0])
	for col := 0; col < numCols; col++ {
		maxHeight := -1
		for row := 0; row < len(trees); row++ {
			currentHeight, err := strconv.Atoi(string(trees[row][col]))
			if err != nil {
				return nil, fmt.Errorf("got error %s, does your trees argument only contain numerical data?", err)
			}
			if currentHeight > maxHeight {
				maxHeight = currentHeight
				visible = append(visible, strconv.Itoa(row)+" "+strconv.Itoa(col))
			}
		}
	}
	return visible, nil
}

// VisibleFromBottom accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from the outside on the bottom of the grid. Each coordinate in the returned
// slice is in the form "R C", where R indicates the row and C indicates the
// column of a visible tree (e.g. "0 3", "1 4", etc.) An error is returned if
// the grid contains any non-numerical data.
func VisibleFromBottom(trees []string) ([]string, error) {
	if len(trees) == 0 {
		return nil, nil
	}
	var visible []string
	numCols := len(trees[0])
	for col := 0; col < numCols; col++ {
		maxHeight := -1
		for row := len(trees) - 1; row >= 0; row-- {
			currentHeight, err := strconv.Atoi(string(trees[row][col]))
			if err != nil {
				return nil, fmt.Errorf("got error %s, does your trees argument only contain numerical data?", err)
			}
			if currentHeight > maxHeight {
				maxHeight = currentHeight
				visible = append(visible, strconv.Itoa(row)+" "+strconv.Itoa(col))
			}
		}
	}
	return visible, nil
}
