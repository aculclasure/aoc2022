package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aculclasure/aoc2022/comms"
)

func main() {
	f, err := os.Open("./day6_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var data string
	scn := bufio.NewScanner(f)
	for scn.Scan() {
		data = scn.Text()
	}
	if err := scn.Err(); err != nil {
		log.Fatal(err)
	}
	marker := comms.StartPacketMarker(data)
	fmt.Printf("The start of marker index is %d\n", marker)
}
