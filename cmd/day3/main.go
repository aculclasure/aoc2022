package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/elf"
)

func main() {
	f, err := os.Open("./day3_input.txt")
	if err != nil {
		log.Fatal(err)
	}

	sum, err := elf.SumDuplicateRucksackItemPriorities(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The sum of priorities of shared items in rucksacks is %d\n", sum)
}
