package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"log"
)

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	n := 0
	inCh := make(chan int, 2)
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			outCh := make(chan int, 1)
			inCh <- x
			inCh <- y
			intcode.Execute(prog, inCh, outCh)
			switch <-outCh {
			case 0:
				fmt.Print(".")
			case 1:
				fmt.Print("#")
				n++
			}
		}
		fmt.Print("\n")
	}
	log.Print(n)
}
