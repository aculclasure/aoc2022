package main

import (
	"fmt"
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
	err = mitm.Run(monkeys, numRounds)
	if err != nil {
		log.Fatal(err)
	}
	business := mitm.MonkeyBusiness(monkeys)
	fmt.Printf("Day 11, Part 1: The level of monkey business after %d rounds is %d\n", numRounds, business)
}
