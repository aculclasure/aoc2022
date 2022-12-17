package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/devices"
)

func main() {
	f, err := os.Open("./day7_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rootDir, err := devices.TreeFromTerminalOutput(f)
	if err != nil {
		log.Fatal(err)
	}
	const maxTotalSizePerDirectory = 100000
	matches := devices.DirectoriesSmallerThan(rootDir, maxTotalSizePerDirectory)
	sum := 0
	for _, m := range matches {
		sum += m.TotalSize()
	}
	fmt.Printf("Sum of total sizes of all directories smaller than %d is %d\n", maxTotalSizePerDirectory, sum)

	const minSystemFreeSpace = 30000000
	bestDirToRemove := rootDir.BestDirectoryToCleanup(minSystemFreeSpace)
	if bestDirToRemove == nil {
		log.Fatal(err)
	}
	fmt.Printf("The best dir to cleanup (%s) has a total size of %d\n", bestDirToRemove.Name, bestDirToRemove.TotalSize())
}
