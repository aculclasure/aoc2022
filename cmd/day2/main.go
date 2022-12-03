package main

import (
	"fmt"
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
}
