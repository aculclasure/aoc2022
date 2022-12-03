package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/rps"
)

func main() {
	f, err := os.Open("./day2_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	score, err := rps.ComputeStrategyScore(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Following the strategy guide results in a score of %d\n", score)

	f.Seek(0, io.SeekStart)
	score, err = rps.ComputeCheatStrategyScore(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Following the cheating strategy guide results in a score of %d\n", score)
}
