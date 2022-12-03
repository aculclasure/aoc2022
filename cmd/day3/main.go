package main

import (
	"fmt"
	"io"
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

	f.Seek(0, io.SeekStart)
	sum, err = elf.SumBadgeItemPriorities(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The sum of badge item priorities for all elf groups is %d\n", sum)
}
