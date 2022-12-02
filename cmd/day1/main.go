package main

import (
	"fmt"
	"io"
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

	maxNumCalories, err := elf.MaxCalories(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The greatest number of calories carried by a single elf is %d\n", maxNumCalories)

	f.Seek(0, io.SeekStart)
	top, err := elf.TopCaloryCarriers(f)
	if err != nil {
		log.Fatal(err)
	}
	sum := 0
	for _, v := range top {
		sum += v
	}
	fmt.Printf("The total number of calories from the top 3 elf carriers is %d\n", sum)
}
