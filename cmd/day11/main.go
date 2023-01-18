package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/mitm"
)

func main() {
	f, err := os.Open("day11_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	monkeys, err := mitm.MonkeysFromInput(f)
	if err != nil {
		log.Fatal(err)
	}
	numRounds := 20
	adjuster1 := mitm.AdjustWorryLevelPart1{Divisor: 3}
	err = mitm.Run(monkeys, numRounds, adjuster1)
	if err != nil {
		log.Fatal(err)
	}
	business := mitm.MonkeyBusiness(monkeys)
	fmt.Printf("Day 11, Part 1: The level of monkey business after %d rounds is %d\n", numRounds, business)

	f.Seek(0, io.SeekStart)
	monkeys, err = mitm.MonkeysFromInput(f)
	if err != nil {
		log.Fatal(err)
	}
	numRounds = 10000
	adjuster2 := mitm.AdjustWorryLevelPart2{CommonMultiple: mitm.CommonMultiple(monkeys)}
	err = mitm.Run(monkeys, numRounds, adjuster2)
	if err != nil {
		log.Fatal(err)
	}
	business = mitm.MonkeyBusiness(monkeys)
	fmt.Printf("Day 11, Part 2: The level of monkey business after %d rounds is %d\n", numRounds, business)
}
