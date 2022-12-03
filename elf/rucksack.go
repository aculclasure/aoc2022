package elf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func FindDuplicateRucksackItems(rucksack string) []rune {
	if len(rucksack) < 2 {
		return nil
	}

	ruckSackItems := []rune(rucksack)
	firstCompartment := ruckSackItems[:len(ruckSackItems)/2]
	secondCompartment := ruckSackItems[len(ruckSackItems)/2:]
	itemSet := map[rune]struct{}{}
	for _, v := range secondCompartment {
		itemSet[v] = struct{}{}
	}

	dups := map[rune]struct{}{}
	for _, v := range firstCompartment {
		if _, ok := itemSet[v]; ok {
			dups[v] = struct{}{}
		}
	}
	var duplicates []rune
	for k := range dups {
		duplicates = append(duplicates, k)
	}

	return duplicates
}

func SumDuplicateRucksackItemPriorities(data io.Reader) (int, error) {
	if data == nil {
		return 0, errors.New("data argument must be non-nil")
	}

	scn := bufio.NewScanner(data)
	priorities := rucksackItemPriorities()
	sum := 0
	for scn.Scan() {
		sharedItems := FindDuplicateRucksackItems(scn.Text())
		for _, v := range sharedItems {
			priorityVal, ok := priorities[v]
			if !ok {
				return 0, fmt.Errorf("shared item %s must have an assigned priority value", string(v))
			}
			sum += priorityVal
		}
	}
	err := scn.Err()
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func rucksackItemPriorities() map[rune]int {
	lowerCaseItems := []rune("abcdefghijklmnopqrstuvwxyz")
	upperCaseItems := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	priorities := make(map[rune]int, len(lowerCaseItems)+len(upperCaseItems))

	for i, v := range lowerCaseItems {
		priorities[v] = i + 1
	}
	for i, v := range upperCaseItems {
		priorities[v] = i + 27
	}

	return priorities
}
