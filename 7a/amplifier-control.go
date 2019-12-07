package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"os"
	"strconv"
	"strings"
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
	log.Printf("trying %v", seq)
	accum := 0
	for _, phase := range seq {
		input := fmt.Sprintf("%d\n%d\n", phase, accum)
		//log.Printf("%q", input)
		output := &bytes.Buffer{}
		if _, err := intcode.Execute(program, bufio.NewReader(bytes.NewBufferString(input)), output); err != nil {
			return -1, fmt.Errorf("during execution of seq %v: %v", seq, err)
		}
		var err error
		accum, err = strconv.Atoi(strings.TrimSpace(output.String()))
		if err != nil {
			return -1, fmt.Errorf("non-integer output from seq %v of %q: %v", seq, output.String(), err)
		}
	}
	return accum, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <amplifier program file>", os.Args[0])
	}

	prog, err := intcode.LoadFile(os.Args[1])
	if err != nil {
		log.Fatalf("loading program: %v", err)
	}

	bestSeq, value, err := try(prog, []int{}, []int{0, 1, 2, 3, 4})
	if err != nil {
		log.Fatalf("executing program: %v", err)
	}
	log.Printf("best phase sequence %v yields output %d", bestSeq, value)
}
