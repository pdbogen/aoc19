package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type OpCode int

const (
	OpCut OpCode = iota
	OpDeal
	OpFlip
)

type Operation struct {
	op OpCode
	n  int
}

func (o Operation) String() string {
	switch o.op {
	case OpCut:
		return fmt.Sprintf("Cut(%d)", o.n)
	case OpDeal:
		return fmt.Sprintf("Deal(%d)", o.n)
	case OpFlip:
		return "Flip()"
	}
	return fmt.Sprintf("UnknownOp(%d)", o.n)
}

func Cut(n int) Operation {
	return Operation{OpCut, n}
}

func Deal(n int) Operation {
	return Operation{OpDeal, n}
}

func Flip() Operation {
	return Operation{op: OpFlip}
}

type Stack struct {
	Cards []int
}

func (s Stack) Cut(n int) Stack {
	newStack := make([]int, len(s.Cards))
	if n >= 0 {
		copy(newStack, s.Cards[n:])
		copy(newStack[len(newStack)-n:], s.Cards)
	} else {
		copy(newStack, s.Cards[len(s.Cards)+n:])
		copy(newStack[-n:], s.Cards)
	}
	return Stack{newStack}
}

func (s Stack) Deal(increment int) Stack {
	newStack := make([]int, len(s.Cards))
	for i := 0; i < len(s.Cards); i++ {
		ti := (i * increment) % len(s.Cards)
		//log.Printf("Deal(%d): %d to %d", increment, Gi, ti)
		newStack[ti] = s.Cards[i]
	}
	return Stack{newStack}
}

func (s Stack) Shuffle(ops []Operation) Stack {
	for _, op := range ops {
		switch op.op {
		case OpCut:
			s = s.Cut(op.n)
		case OpDeal:
			s = s.Deal(op.n)
		case OpFlip:
			s = s.Flip()
		}
		log.Printf("after %s: %v", op, s)
	}
	return s
}

func (s Stack) Flip() Stack {
	newStack := make([]int, len(s.Cards))
	for i := 0; i < len(s.Cards); i++ {
		newStack[i] = s.Cards[len(s.Cards)-1-i]
	}
	return Stack{newStack}
}

func LoadShuffle(filename string) ([]Operation, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var ret []Operation
	lines := strings.Split(string(f), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "deal with increment ") {
			n, err := strconv.Atoi(strings.TrimPrefix(line, "deal with increment "))
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %v", line, err)
			}
			ret = append(ret, Deal(n))
		} else if strings.HasPrefix(line, "cut ") {
			n, err := strconv.Atoi(strings.TrimPrefix(line, "cut "))
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %v", line, err)
			}
			ret = append(ret, Cut(n))
		} else if line == "deal into new stack" {
			ret = append(ret, Flip())
		}
	}
	return ret, nil
}

func main() {
	var stack Stack
	for i := 0; i <= 10006; i++ {
		stack.Cards = append(stack.Cards, i)
	}
	log.Print(len(stack.Cards))
	ops, err := LoadShuffle("input")
	if err != nil {
		log.Fatal(err)
	}
	stack = stack.Shuffle(ops)
	log.Print(stack.Cards)
	for i, c := range stack.Cards {
		if c == 2019 {
			log.Printf("2019 @ card %d", i)
			break
		}
	}
}
