package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"strings"
)

// tt
//
// A  B  C  D  J
// T  T  T  T  F
// T  T  T  F  F
// T  T  F  T  T
//
//

var program = strings.Join([]string{
	"NOT A T", // jump if A empty
	"AND D T", // only if there's a place to land
	"OR T J", // jump in either case
	"NOT B T", // jump if B empty
	"AND D T", // only if there's a place to land
	"OR T J", // jump in either case
	"NOT C T", // jump if B empty
	"AND D T", // only if there's a place to land
	"OR T J", // jump in either case
	"WALK",
	""}, "\n")

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	inCh := make(chan int)
	outCh := make(chan int)
	go intcode.Execute(prog, inCh, outCh)
	done := make(chan struct{})
	go func() {
		for i := range outCh {
			if i < 128 {
				print(string(byte(i)))
			} else {
				log.Print(i)
			}
		}
		close(done)
	}()

	for _, b := range []byte(program) {
		inCh <- int(b)
	}
	<-done
}
