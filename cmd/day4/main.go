package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/camp"
)

func main() {
	f, err := os.Open("./day4_input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pairs, err := camp.GetFullyOverlappingPairs(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The number of fully overlapping pairs is %d\n", len(pairs))

	f.Seek(0, io.SeekStart)
	pairs, err = camp.GetOverlappingPairs(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The number of overlapping pairs is %d\n", len(pairs))
}
