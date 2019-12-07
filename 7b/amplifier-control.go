package main

import (
	"errors"
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"os"
)

func try(program []int, seq []int, rem []int) (best []int, result int, err error) {
	if len(rem) == 0 {
		value, err := runTry(program, seq)
		if err != nil {
			return nil, 0, err
		}
		return seq, value, nil
	}

	var bestSeq []int
	var bestValue int
	for i, next := range rem {
		subSeq := make([]int, len(seq)+1)
		copy(subSeq, seq)
		subSeq[len(seq)] = next

		subRem := make([]int, len(rem)-1)
		copy(subRem, rem[:i])
		copy(subRem[i:], rem[i+1:])

		thisSeq, thisValue, err := try(program, subSeq, subRem)
		if err != nil {
			return nil, 0, err
		}
		if thisValue > bestValue {
			bestSeq = thisSeq
			bestValue = thisValue
		}
	}
	return bestSeq, bestValue, nil
}

func runTry(program []int, seq []int) (int, error) {
	if len(seq) < 2 {
		return -1, errors.New("cannot create feedback loop without at least two amplifiers")
	}

	log.Printf("trying %v", seq)
	errCh := make(chan error)
	result := make(chan int)
	var inputs []chan int
	for _, phase := range seq {
		inCh := make(chan int, 2)
		inCh <- phase
		inputs = append(inputs, inCh)
	}
	for i := range seq {
		if i == 0 {
			inputs[i] <- 0
			go func() {
				if _, err := intcode.Execute(program, inputs[0], inputs[1]); err != nil {
					errCh <- err
				} else {
					result <- <-inputs[0]
				}
			}()
		} else {
			go func(i int) {
				if _, err := intcode.Execute(program, inputs[i], inputs[(i+1)%len(seq)]); err != nil {
					errCh <- err
				}
			}(i)
		}
	}

	select {
	case yield := <-result:
		log.Printf("yielded %d", yield)
		return yield, nil
	case err := <-errCh:
		return -1, err
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <amplifier program file>", os.Args[0])
	}

	prog, err := intcode.LoadFile(os.Args[1])
	if err != nil {
		log.Fatalf("loading program: %v", err)
	}

	bestSeq, value, err := try(prog, []int{}, []int{5, 6, 7, 8, 9})
	if err != nil {
		log.Fatalf("executing program: %v", err)
	}
	log.Printf("best phase sequence %v yields output %d", bestSeq, value)
}
