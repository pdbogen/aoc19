package intcode

import (
	"bytes"
	"fmt"
	. "github.com/pdbogen/aoc19/intcode/opcodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Computer struct {
	RelativeBase int
	Program      []int
}
type Op struct {
	Code  int
	Modes []int
	Args  []int

	// When an argument represents an address, this is the canonical address
	Address []int

	// This is the interpretation of the argument as a value
	BakedArgs []int
}

func (c *Computer) ReadOp(ptr int, args int) (op *Op, err error) {
	if ptr+args >= len(c.Program) {
		return nil, fmt.Errorf("expected %d arguments but not enough program left at %d", args, ptr)
	}

	ret := &Op{
		Code:  c.Program[ptr] % 100,
		Modes: []int{},
	}

	flags := c.Program[ptr] / 100
	for i := 0; i < args; i++ {
		ret.Args = append(ret.Args, c.Program[ptr+1+i])
		ret.Modes = append(ret.Modes, flags%10)
		flags = flags / 10
	}

	ret.BakedArgs = make([]int, len(ret.Args))
	ret.Address = make([]int, len(ret.Args))
	for argN, arg := range ret.Args {
		switch ret.Modes[argN] {
		case ModePosition:
			if arg < len(c.Program) {
				ret.BakedArgs[argN] = c.Program[arg]
			}
			ret.Address[argN] = arg
		case ModeImmediate:
			ret.BakedArgs[argN] = arg
			ret.Address[argN] = ptr + argN + 1
		case ModeRelative:
			if c.RelativeBase+arg < len(c.Program) {
				ret.BakedArgs[argN] = c.Program[c.RelativeBase+arg]
			}
			ret.Address[argN] = c.RelativeBase + arg
		default:
			return nil, fmt.Errorf("unhandled parameter mode %d", ret.Modes[argN])
		}
	}

	return ret, nil
}

func (c Computer) Jump(ptr int, wantNonZero bool) (out []int, outPtr int, err error) {
	op, err := c.ReadOp(ptr, 2)
	if err != nil {
		return nil, 0, err
	}

	isNonZero := op.BakedArgs[0] != 0
	if isNonZero == wantNonZero {
		return c.Program, op.BakedArgs[1], nil
	}
	return c.Program, ptr + 3, nil
}

// A compare argument takes three arguments; two to compare, one to save the result.
func (c Computer) Compare(ptr int, comparison func(i, j int) bool) (out []int, outPtr int, err error) {
	op, err := c.ReadOp(ptr, 3)
	if err != nil {
		return nil, 0, err
	}

	result := 0
	if comparison(op.BakedArgs[0], op.BakedArgs[1]) {
		result = 1
	}

	if op.Address[2] >= len(c.Program) {
		out = make([]int, op.Address[2]+1)
	} else {
		out = make([]int, len(c.Program))
	}
	copy(out, c.Program)
	out[op.Address[2]] = result

	return out, ptr + 4, nil
}

// Input takes one argument, the address where we'll save the read-in value.
func (c Computer) Input(ptr int, inChan <-chan int) (out []int, outPtr int, err error) {
	op, err := c.ReadOp(ptr, 1)
	if err != nil {
		return nil, 0, err
	}

	inputInt := <-inChan
	dest := op.Address[0]
	outPtr = ptr + 2

	if dest+1 >= len(c.Program) {
		out = make([]int, dest+1)
	} else {
		out = make([]int, len(c.Program))
	}

	copy(out, c.Program)
	out[dest] = inputInt

	return out, outPtr, nil
}

// Output takes one argument, the value to print out.
func (c Computer) Output(ptr int, outChan chan<- int) (out []int, outPtr int, err error) {
	op, err := c.ReadOp(ptr, 1)
	if err != nil {
		return nil, 0, err
	}
	if op.Code != OpOutput {
		return nil, 0, fmt.Errorf("read op code %d is not %d", op.Code, OpOutput)
	}

	if outChan == nil {
		log.Printf("output: %d", op.BakedArgs[0])
	} else {
		outChan <- op.BakedArgs[0]
	}

	return c.Program, ptr + 2, nil
}

func (c Computer) BinaryMath(opCode int, ptr int, op func(int, int) int) (out []int, outPtr int, err error) {
	readOp, err := c.ReadOp(ptr, 3)
	if err != nil {
		return nil, 0, err
	}

	if readOp.Code != opCode {
		return nil, 0, fmt.Errorf("expected opcode %d at ptr %d but found %d", opCode, ptr, readOp.Code)
	}

	dest := readOp.Address[2]
	outPtr = ptr + 4

	if dest >= len(c.Program) {
		out = make([]int, dest+1)
	} else {
		out = make([]int, len(c.Program))
	}
	copy(out, c.Program)
	out[dest] = op(readOp.BakedArgs[0], readOp.BakedArgs[1])

	return out, outPtr, nil
}

func (c *Computer) SetRelativeBase(ptr int) (out []int, outPtr int, err error) {
	op, err := c.ReadOp(ptr, 1)
	if err != nil {
		return nil, 0, err
	}

	c.RelativeBase += op.BakedArgs[0]

	return c.Program, ptr + 2, nil
}

func Execute(in []int, input <-chan int, output chan<- int) (out []int, err error) {
	c := Computer{Program: in}

	if output != nil {
		defer close(output)
	}

	ptr := 0
execution:
	for ptr < len(c.Program) {
		// log.Printf("%d in %v [%d] %v", ptr, prog[:ptr], prog[ptr], prog[ptr+1:])
		var newPtr int
		var newProg []int
		switch c.Program[ptr] % 100 {
		case OpAdd:
			newProg, newPtr, err = c.BinaryMath(OpAdd, ptr, func(a int, b int) int { return a + b })
		case OpMul:
			newProg, newPtr, err = c.BinaryMath(OpMul, ptr, func(a int, b int) int { return a * b })
		case OpInput:
			newProg, newPtr, err = c.Input(ptr, input)
		case OpOutput:
			newProg, newPtr, err = c.Output(ptr, output)
		case OpJumpTrue:
			newProg, newPtr, err = c.Jump(ptr, true)
		case OpJumpFalse:
			newProg, newPtr, err = c.Jump(ptr, false)
		case OpLessThan:
			newProg, newPtr, err = c.Compare(ptr, func(i, j int) bool { return i < j })
		case OpEquals:
			newProg, newPtr, err = c.Compare(ptr, func(i, j int) bool { return i == j })
		case OpSetRelativeBase:
			newProg, newPtr, err = c.SetRelativeBase(ptr)
		case OpEnd:
			break execution
		default:
			return nil, fmt.Errorf("unrecognized opcode %d at %d", c.Program[ptr], ptr)
		}
		if err != nil {
			return nil, fmt.Errorf("at ptr %d in %v, opcode %d resulted in error: %v", ptr, c.Program, c.Program[ptr], err)
		}
		c.Program = newProg
		ptr = newPtr
	}
	return c.Program, nil
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

	programText := strings.Split(strings.TrimSpace(string(prog)), ",")
	program := make([]int, len(programText))
	for i, s := range programText {
		if s == "" {
			continue
		}
		program[i], err = strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("token %d %q could not be converted to integer: %v", i, s, err)
		}
	}

	return program, nil
}

func Provider(values []int) <-chan int {
	ret := make(chan int)
	go func(ch chan<- int, values []int) {
		defer close(ret)
		for _, v := range values {
			ret <- v
		}
	}(ret, values)
	return ret
}
