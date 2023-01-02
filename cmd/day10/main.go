package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/devices"
)

func main() {
	f, err := os.Open("day10_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sigStrengths, err := devices.SignalStrengths(f)
	if err != nil {
		log.Fatal(err)
	}
	sum := 0
	for _, s := range sigStrengths {
		sum += s
	}
	fmt.Printf("Part 1: The sum of signal strengths is %d\n", sum)
	f.Seek(0, io.SeekStart)
	output, err := devices.DrawOnScreen(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 2: Output")
	fmt.Println(output)
}
