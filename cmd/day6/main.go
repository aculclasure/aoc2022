package main

import (
	"bufio"
	"fmt"
	"io"
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
	fmt.Printf("The start of packet marker index is %d\n", marker)

	f.Seek(0, io.SeekStart)
	scn = bufio.NewScanner(f)
	for scn.Scan() {
		data = scn.Text()
	}
	if err := scn.Err(); err != nil {
		log.Fatal(err)
	}
	marker = comms.StartMessageMarker(data)
	fmt.Printf("The start of message marker index is %d\n", marker)
}
