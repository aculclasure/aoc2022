package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/rope"
)

func main() {
	f, err := os.Open("day9_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r, err := rope.Run(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The number of positions visited by the rope's tail end is %d\n", len(r.Tail.Visited))
}
