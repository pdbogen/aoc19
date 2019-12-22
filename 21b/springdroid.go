package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"strings"
)

// tt
//
// A  B  C  D  E F G H I J
// T  T  F  T  F       T T
//

var program = strings.Join([]string{
	"NOT A T",
	"OR  T J",
	"NOT B T",
	"OR  T J",
	"NOT C T",
	"OR  T J",
	"AND D J",
	"AND E T",
	"OR  H T",
	"AND T J",
	"RUN",
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
