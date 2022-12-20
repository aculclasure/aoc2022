package camp

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// result represents the result of an execution of the VisibleFromLeft(),
// VisibleFromRight(), VisibleFromTop(), or VisibleFromBottom() functions when it
// is invoked as a goroutine. If the call is successful, then the coordinates of
// the visible trees are stored in the visTrees field. If the call is unsuccessful,
// then the error is stored in the err field.
type result struct {
	visTrees []string
	err      error
}

// TreesFromBytes accepts a slice of bytes that is intended to come from a file
// containing line-separated tree height data and returns a slice where each string
// is a row of tree height data.
func TreesFromBytes(input []byte) []string {
	return strings.Fields(strings.TrimSpace(string(input)))
}

// AllVisibleTrees accepts a slice of strings representing a grid of tree height
// data and returns a slice of coordinate strings for all trees that are visible
// from outside the grid. Each coordinate in the returned slice is in the form
// "R C", where R indicates the row and C indicates the column of a visible tree
// (e.g. "0 3", "1 4", etc.) An error is returned if the grid contains any
// non-numerical data.
func AllVisibleTrees(trees []string) ([]string, error) {
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
