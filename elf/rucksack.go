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

func FindBadgeInGroup(group [][]rune) (rune, error) {
	if len(group) < 2 {
		return '0', fmt.Errorf("group size must be at least 2 (got %d)", len(group))
	}

	type set map[rune]struct{}
	var sets []set
	for _, grp := range group {
		itemSet := set{}
		for _, item := range grp {
			itemSet[item] = struct{}{}
		}
		sets = append(sets, itemSet)
	}

	firstSet := sets[0]
	isBadgeItem := false
	for item := range firstSet {
		for _, nextSet := range sets[1:] {
			if _, ok := nextSet[item]; !ok {
				isBadgeItem = false
				break
			}
			isBadgeItem = true
		}
		if isBadgeItem {
			return item, nil
		}
	}

	return '0', errors.New("unable to locate a badge item type in the given group")
}

func SumBadgeItemPriorities(data io.Reader) (int, error) {
	if data == nil {
		return 0, errors.New("data argument must be non-nil")
	}

	const groupSize = 3
	var (
		group        [][]rune
		numLinesRead int
		sum          int
		scn          = bufio.NewScanner(data)
		priorities   = rucksackItemPriorities()
	)
	for scn.Scan() {
		ruckSackItems := []rune(scn.Text())
		group = append(group, ruckSackItems)
		numLinesRead++
		if numLinesRead%groupSize == 0 {
			badge, err := FindBadgeInGroup(group)
			if err != nil {
				return 0, err
			}
			badgeVal, ok := priorities[badge]
			if !ok {
				return 0, fmt.Errorf("badge item %s must have an assigned priority value", string(badge))
			}
			sum += badgeVal
			group = [][]rune{}
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
