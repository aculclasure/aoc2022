package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/camp"
)

func main() {
	f, err := os.Open("./day4_input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fullOverlaps, err := camp.GetFullyOverlappingPairs(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The number of fully overlapping pairs is %d\n", len(fullOverlaps))
}
