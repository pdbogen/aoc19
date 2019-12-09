package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"sync"
)

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatalf("could not load input: %v", err)
	}

	outCh := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done();
		for i := range outCh {
			log.Print(i)
		}
	}()

	inCh := make(chan int, 1)
	inCh <- 1 // test mode
	close(inCh)

	_, err = intcode.Execute(prog, inCh, outCh)
	wg.Wait()
	if err != nil {
		log.Fatalf("error during execution: %v", err)
	}
}
