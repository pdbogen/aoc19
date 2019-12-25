package main

import (
	"bufio"
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"os"
	"strings"
)

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	inCh := make(chan int)
	outCh := make(chan int)
	go func() {
		for out := range outCh {
			fmt.Print(string(byte(out)))
		}
	}()
	go intcode.Execute(prog, inCh, outCh)
	buf := bufio.NewReader(os.Stdin)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		for _, c := range line {
			inCh <- int(c)
		}
		inCh <- 10
	}
}
