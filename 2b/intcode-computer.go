package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	OpAdd = 1
	OpMul = 2
	OpEnd = 99
)

func Add(chain []int, ptr int) (resultChain []int, resultPtr int, err error) {
	return binOp(OpAdd, func(a, b int) int { return a + b })(chain, ptr)
}

func Mul(chain []int, ptr int) (resultChain []int, resultPtr int, err error) {
	return binOp(OpMul, func(a, b int) int { return a * b })(chain, ptr)
}

func binOp(opCode int, op func(int, int) int) func(in []int, ptr int) (out []int, outPtr int, err error) {
	return func(in []int, ptr int) (out []int, outPtr int, err error) {
		if ptr+3 >= len(in) {
			return nil, 0, fmt.Errorf("program pointer %d out of range", ptr)
		}

		if in[ptr] != opCode {
			return nil, 0, fmt.Errorf("program pointed %d pointed to %d, not opcode %d", ptr, in[ptr], opCode)
		}

		opAIdx := in[ptr+1]
		if opAIdx > len(in) {
			return nil, 0, fmt.Errorf("op A index %d out of range", opAIdx)
		}
		opA := in[opAIdx]

		opBIdx := in[ptr+2]
		if opBIdx > len(in) {
			return nil, 0, fmt.Errorf("op B index %d out of range", opBIdx)
		}
		opB := in[opBIdx]

		dest := in[ptr+3]
		outPtr = ptr + 4

		if dest+1 >= len(in) {
			out = make([]int, dest+1)
		} else {
			out = make([]int, len(in))
		}

		copy(out, in)

		out[dest] = op(opA, opB)

    // log.Printf("wrote %d (@%d) op%d %d (@%d) = %d (@%d)", opA, opAIdx, opCode, opB, opBIdx, out[dest], dest)

		return out, outPtr, nil
	}
}

func execute(in []int) (out []int, err error) {
	ptr := 0
	prog := in
execution:
	for ptr < len(prog) {
		// log.Printf("%d in %v [%d] %v", ptr, prog[:ptr], prog[ptr], prog[ptr+1:])
		var newPtr int
		var newProg []int
		switch prog[ptr] {
		case OpAdd:
			newProg, newPtr, err = Add(prog, ptr)
		case OpMul:
			newProg, newPtr, err = Mul(prog, ptr)
		case OpEnd:
			break execution
		}
		if err != nil {
			return nil, fmt.Errorf("at ptr %d in %v, opcode %d resulted in error %v", ptr, prog, prog[ptr], err)
		}
		prog = newProg
		ptr = newPtr
	}
	return prog, nil
}

func main() {
	prog, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read program on stdin: %v", err)
	}
	chainStr := strings.Split(strings.TrimSpace(string(prog)), ",")
	chainInt := make([]int, len(chainStr))
	for i, s := range chainStr {
		if s == "" {
			continue
		}
		chainInt[i], err = strconv.Atoi(s)
		if err != nil {
			log.Fatalf("token %d %q could not be converted to integer: %v", i, s, err)
		}
	}

	// log.Printf("going to execute %v", chainInt)
	result, err := execute(chainInt)
	if err != nil {
		log.Fatalf("error during execution: %+v", err)
	}
	// log.Printf("got back %v", result)

	first := true
	for _, i := range result {
		if !first {
			fmt.Printf(",%d", i)
		} else {
			fmt.Printf("%d", i)
			first = false
		}
	}
	fmt.Print("\n")
}
