package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/elf"
)

func main() {
	f, err := os.Open("./day1_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	topCounts, err := elf.TopCaloryCounts(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The greatest number of calories carried by a single elf is %d\n", topCounts[0])

	sum := 0
	for _, v := range topCounts {
		sum += v
	}
	fmt.Printf("The total number of calories from the top 3 elf carriers is %d\n", sum)
}
