package intcode

import (
	"bytes"
	"fmt"
	. "github.com/pdbogen/aoc19/intcode/opcodes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Op struct {
	Code      int
	Modes     []int
	Args      []int
	BakedArgs []int
}

func ReadOp(in []int, ptr int, args int) (op *Op, err error) {
	if ptr+args >= len(in) {
		return nil, fmt.Errorf("expected %d arguments but not enough program left at %d", args, ptr)
	}

	ret := &Op{
		Code:  in[ptr] % 100,
		Modes: []int{},
	}

	flags := in[ptr] / 100
	for i := 0; i < args; i++ {
		ret.Args = append(ret.Args, in[ptr+1+i])
		ret.Modes = append(ret.Modes, flags%10)
		flags = flags / 10
	}

	ret.BakedArgs = make([]int, len(ret.Args))
	for argN, arg := range ret.Args {
		ret.BakedArgs[argN] = arg
		if ret.Modes[argN] == ModePosition {
			if arg < len(in) {
				ret.BakedArgs[argN] = in[arg]
			}
		} else if ret.Modes[argN] != ModeImmediate {
			return nil, fmt.Errorf("unhandled parameter mode %d", ret.Modes[argN])
		}
	}

	return ret, nil
}

func Jump(in []int, ptr int, wantNonZero bool) (out []int, outPtr int, err error) {
	op, err := ReadOp(in, ptr, 2)
	if err != nil {
		return nil, 0, err
	}

	isNonZero := op.BakedArgs[0] != 0
	if isNonZero == wantNonZero {
		return in, op.BakedArgs[1], nil
	}
	return in, ptr + 3, nil
}

func Compare(in []int, ptr int, comparison func(i, j int) bool) (out []int, outPtr int, err error) {
	op, err := ReadOp(in, ptr, 3)
	if err != nil {
		return nil, 0, err
	}

	result := 0
	if comparison(op.BakedArgs[0], op.BakedArgs[1]) {
		result = 1
	}

	if op.Args[2] >= len(in) {
		out = make([]int, op.Args[2]+1)
	} else {
		out = make([]int, len(in))
	}
	copy(out, in)
	out[op.Args[2]] = result

	return out, ptr + 4, nil
}

func Input(in []int, ptr int, inChan <-chan int) (out []int, outPtr int, err error) {
	if ptr+1 >= len(in) {
		return nil, 0, fmt.Errorf("program pointer %d out of range", ptr)
	}
	if in[ptr] != OpInput {
		return nil, 0, fmt.Errorf("opcode %d at ptr %d is not %d", in[ptr], ptr, OpInput)
	}

	inputInt := <-inChan

	dest := in[ptr+1]
	outPtr = ptr + 2

	if dest+1 >= len(in) {
		out = make([]int, dest+1)
	} else {
		out = make([]int, len(in))
	}

	copy(out, in)
	out[dest] = inputInt

	return out, outPtr, nil
}

func Output(in []int, ptr int, outChan chan<- int) (out []int, outPtr int, err error) {
	op, err := ReadOp(in, ptr, 1)
	if err != nil {
		return nil, 0, err
	}
	if op.Code != OpOutput {
		return nil, 0, fmt.Errorf("read op code %d is not %d", op.Code, OpOutput)
	}

	if op.Args[0] >= len(in) {
		return nil, 0, fmt.Errorf("output target %d out of range (chain length %d)", op.BakedArgs[0], len(in))
	}

	outChan <- op.BakedArgs[0]

	return in, ptr + 2, nil
}

func Add(chain []int, ptr int) (resultChain []int, resultPtr int, err error) {
	return binOp(OpAdd, func(a, b int) int { return a + b })(chain, ptr)
}

func Multiply(chain []int, ptr int) (resultChain []int, resultPtr int, err error) {
	return binOp(OpMul, func(a, b int) int { return a * b })(chain, ptr)
}

func binOp(opCode int, op func(int, int) int) func(in []int, ptr int) (out []int, outPtr int, err error) {
	return func(in []int, ptr int) (out []int, outPtr int, err error) {
		readOp, err := ReadOp(in, ptr, 3)
		if err != nil {
			return nil, 0, err
		}
		if readOp.Code != opCode {
			return nil, 0, fmt.Errorf("program pointed %d pointed to %d, not opcode %d", ptr, readOp.Code, opCode)
		}

		opAIdx := readOp.Args[0]
		if readOp.Modes[0] == ModeImmediate {
			opAIdx = ptr + 1
		}
		if opAIdx >= len(in) {
			return nil, 0, fmt.Errorf("op A index %d out of range", opAIdx)
		}
		opA := in[opAIdx]

		opBIdx := readOp.Args[1]
		if readOp.Modes[1] == ModeImmediate {
			opBIdx = ptr + 2
		}
		if opBIdx >= len(in) {
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

func Execute(in []int, input <-chan int, output chan<- int) (out []int, err error) {
	defer close(output)
	ptr := 0
	prog := in
execution:
	for ptr < len(prog) {
		// log.Printf("%d in %v [%d] %v", ptr, prog[:ptr], prog[ptr], prog[ptr+1:])
		var newPtr int
		var newProg []int
		switch prog[ptr] % 100 {
		case OpAdd:
			newProg, newPtr, err = Add(prog, ptr)
		case OpMul:
			newProg, newPtr, err = Multiply(prog, ptr)
		case OpInput:
			newProg, newPtr, err = Input(prog, ptr, input)
		case OpOutput:
			newProg, newPtr, err = Output(prog, ptr, output)
		case OpJumpTrue:
			newProg, newPtr, err = Jump(prog, ptr, true)
		case OpJumpFalse:
			newProg, newPtr, err = Jump(prog, ptr, false)
		case OpLessThan:
			newProg, newPtr, err = Compare(prog, ptr, func(i, j int) bool { return i < j })
		case OpEquals:
			newProg, newPtr, err = Compare(prog, ptr, func(i, j int) bool { return i == j })
		case OpEnd:
			break execution
		default:
			return nil, fmt.Errorf("unrecognized opcode %d at %d", prog[ptr], ptr)
		}
		if err != nil {
			return nil, fmt.Errorf("at ptr %d in %v, opcode %d resulted in error: %v", ptr, prog, prog[ptr], err)
		}
		prog = newProg
		ptr = newPtr
	}
	return prog, nil
}

func MustLoadString(in string) []int {
	out, err := LoadString(in)
	if err != nil {
		panic(err)
	}
	return out
}

func LoadString(in string) (out []int, err error) {
	return Load(bytes.NewBufferString(in))
}

func LoadFile(filename string) (out []int, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening %q for reading: %v", filename, err)
	}
	defer f.Close()
	return Load(f)
}

func Load(in io.Reader) (out []int, err error) {
	prog, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("could not read program: %v", err)
	}

	chainStr := strings.Split(strings.TrimSpace(string(prog)), ",")
	chainInt := make([]int, len(chainStr))
	for i, s := range chainStr {
		if s == "" {
			continue
		}
		chainInt[i], err = strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("token %d %q could not be converted to integer: %v", i, s, err)
		}
	}

	return chainInt, nil
}

func Provider(values []int) (<-chan int) {
	ret := make(chan int)
	go func(ch chan<- int, values []int) {
		defer close(ret)
		for _, v := range values {
			ret <- v
		}
	}(ret, values)
	return ret
}
