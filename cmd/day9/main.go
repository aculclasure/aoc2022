package main

import (
	"fmt"
	"io"
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
	r, err := rope.Run(f, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1: The number of positions visited by the rope's tail end when the rope has 2 knots is %d\n", len(r.Tail.Visited))
	f.Seek(0, io.SeekStart)
	r, err = rope.Run(f, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 2: The number of positions visited by the rope's tail end when the rope has 10 knots is %d\n", len(r.Tail.Visited))
}
