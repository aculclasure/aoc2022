package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/camp"
)

func main() {
	treeData, err := os.ReadFile("day8_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	trees := camp.TreesFromBytes(treeData)
	vis, err := camp.AllVisibleTrees(trees)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The number of trees visible from outside the grid is %d\n", len(vis))
	maxScenicScore, err := camp.MaxScenicScore(trees)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The max scenic score for a tree in the grid is %d\n", maxScenicScore)
}
