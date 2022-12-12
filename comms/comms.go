package comms

import (
	"errors"
)

// HasUniqueChars accepts a slice of runes representing a data stream in the
// elf's communications device and returns true if all the characters in the
// input are unique and returns false otherwise. An error is returned if the
// input is nil or empty.
func HasUniqueChars(input []rune) (bool, error) {
	if len(input) == 0 {
		return false, errors.New("input must be non-empty")
	}

	if len(input) < 2 {
		return true, nil
	}

	set := map[rune]struct{}{}
	for _, r := range input {
		set[r] = struct{}{}
	}
	return len(set) == len(input), nil
}

// StartPacketMarker accepts a string representing a data stream in the elf's
// communications device and returns a number indicating the index of the first
// start of marker packet position. A negative value is returned if no start of
// marker packet can be found in the input.
func StartPacketMarker(input string) int {
	const minNumUniqueChars int = 4
	runes := []rune(input)
	if len(runes) < minNumUniqueChars {
		return -1
	}
	for i := minNumUniqueChars; i <= len(runes); i++ {
		chunk := runes[i-minNumUniqueChars : i]
		unique, err := HasUniqueChars(chunk)
		switch {
		case err != nil:
			fallthrough
		case !unique:
			continue
		default:
			return i
		}
	}
	return -1
}
