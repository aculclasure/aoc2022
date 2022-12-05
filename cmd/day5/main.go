package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/cargo"
)

func main() {
	f, err := os.Open("./day5_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	layout, err := cargo.LayoutFromData(f)
	if err != nil {
		log.Fatal(err)
	}
	topItems := layout.GetTopItems()
	fmt.Printf("The top items from all stacks are %s\n", topItems)
}
