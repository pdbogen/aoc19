package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
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

// Cut 1  -> 1 2 3 4 5 6 7 8 9 0
// Flip   -> 0 9 8 7 6 5 4 3 2 1

// Flip   -> 9 8 7 6 5 4 3 2 1 0
// Cut -1 -> 0 9 8 7 6 5 4 3 2 1

// Cut 1  -> 1 2 3 4 5 6 7 8 9 0
// Cut 1  -> 2 3 4 5 6 7 8 9 0 1

//-1: 9 0 1 2 3 4 5 6 7 8
// 0: 0 1 2 3 4 5 6 7 8 9
// 1: 1 2 3 4 5 6 7 8 9 0
// DoCut returns the index where the card at `idx` will end up after `Cut(n)` with a deck size of `deckSize`
func DoCut(n, idx, deckSize int) int {
	if n < 0 {
		return DoCut(deckSize+n, idx, deckSize)
	}
	return (idx + deckSize - n) % deckSize
}

func UndoCut(n, idx, deckSize int) int {
	return DoCut(-n, idx, deckSize)
}

// Deal 3  -> 0 7 4 1 8 5 2 9 6 3
// Flip    -> 3 6 9 2 5 8 1 4 7 0

// Flip    -> 9 8 7 6 5 4 3 2 1 0
// Deal  3 -> 9 2 5 8 1 4 7 0 3 6
// Cut  -2 -> 3 6 9 2 5 8 1 4 7 0

// Deal 7  -> 0 3 6 9 2 5 8 1 4 7
// Flip    -> 7 4 1 8 5 2 9 6 3 0

// Flip    -> 9 8 7 6 5 4 3 2 1 0
// Deal 7  -> 9 6 3 0 7 4 1 8 5 2
// Cut  -6 -> 7 4 1 8 5 2 9 6 3 0

// Deal 7  -> 0 3 6 9 2 5 8 1 4 7
// Cut  1  -> 3 6 9 2 5 8 1 4 7 0

// Cut  3  -> 3 4 5 6 7 8 9 0 1 2
// Deal 7  -> 3 6 9 2 5 8 1 4 7 0

// Deal 7  -> 0 3 6 9 2 5 8 1 4 7
// Cut  2  -> 6 9 2 5 8 1 4 7 0 3

// Cut  6  -> 6 7 8 9 0 1 2 3 4 5
// Deal 7  -> 6 9 2 5 8 1 4 7 0 3

// Deal 7  -> 0 3 6 9 2 5 8 1 4 7
// Cut  3  -> 9 2 5 8 1 4 7 0 3 6

// 7,1 -> 3,7
// 7,2 -> 6,7
// 7,3 -> 9,7
// 3,1 -> 7,3
// 3,2 -> 4,3 (14%10)
// Deal(n),Cut(m) ->  Cut((deckSize-n)*m), Deal(n)
//(deckSize-m)*n

// Deal n  -> x*n % deckSize
// Cut m   -> (x*n-m+deckSize) % deckSize

// Cut n   -> x-m+deckSize
// Deal m  -> (x-m+DeckSize)*n

// 999-1+10 = 998+10 == 1008 % 10 == 8
// 9-1+10   = 8+10   == 18   % 10 == 8

// Deal 3  -> 0 7 4 1 8 5 2 9 6 3
// Cut  1  -> 7 4 1 8 5 2 9 6 3 0
// Cut  7  -> 7 8 9 0 1 2 3 4 5 6
// Deal 3  -> 7 4 1 8 5 2 9 6 3 0

// Deal 3  -> 0 7 4 1 8 5 2 9 6 3
// Cut  2  -> 4 1 8 5 2 9 6 3 0 7

// Cut  4  -> 4 5 6 7 8 9 0 1 2 3
// Deal 3  -> 4 1 8 5 2 9 6 3 0 7

// Deal n  -> x*n
// Deal m   -> x*n*m
// Deal n*m

func DoDeal(n, idx, deckSize int) int {
	return (idx * n) % deckSize
}

// Deal 3  -> 0 7 4 1 8 5 2 9 6 3
// UndoDeal(3, 0, 10) == 0
// UndoDeal(3, 1, 10) == 7 (7*1 % 10 == 7)
// UndoDeal(3, 3, 10) == 1 (7*3 % 10 == 1)

// inverse returns the value `p` such that, for x mod n == x, p*x mod n == 1
func inverse(x, mod int) *big.Int {
	if x == 1 {
		return big.NewInt(1)
	}

	q := []int{}
	r := []int{}
	p := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
	}
	i := -1
	a := mod
	b := x
	for b != 0 {
		i++
		q = append(q, a/b)
		r = append(r, a%b)
		a, b = b, a%b
		if i >= 2 {
			pi0 := (&big.Int{}).Set(p[i-2])
			pi1 := (&big.Int{}).Set(p[i-1])
			pi1.Mul(pi1, big.NewInt(int64(q[i-2])))
			pi0.Sub(pi0, pi1)
			for pi0.Cmp(big.NewInt(0)) < 0 {
				pi0.Add(pi0, big.NewInt(int64(mod)))
			}
			p = append(p, pi0)
		}
	}

	pi0 := (&big.Int{}).Set(p[i-1])
	pi1 := (&big.Int{}).Set(p[i])
	pi1.Mul(pi1, big.NewInt(int64(q[i-1])))
	pi0.Sub(pi0, pi1)
	for pi0.Cmp(big.NewInt(0)) < 0 {
		pi0.Add(pi0, big.NewInt(int64(mod)))
	}
	p = append(p, pi0) // r[i], 0
	// r[i-1], last non-zero remainder
	// p[i+1], inverse

	if r[i-1] != 1 {
		panic(fmt.Sprintf("cannot compute modular inverse of %d mod %d; gcd is %d, not 1", x, mod, r[i]))
	}

	if !p[i+1].IsInt64() {
		panic("in inverse, " + p[i+1].String() + " overflows int64")
	}
	return p[i+1]
}

func UndoDeal(n, idx, deckSize int) int {
	if idx == 0 {
		return 0
	}

	for n < 0 {
		n += deckSize
	}
	inv := inverse(n, deckSize)
	inv.Mul(inv, big.NewInt(int64(idx)))
	inv.Mod(inv, big.NewInt(int64(deckSize)))
	if !inv.IsInt64() {
		panic("in UndoDeal, " + inv.String() + " overflows int64")
	}
	return int(inv.Int64())
}

func DoFlip(idx, deckSize int) int {
	return deckSize - 1 - idx
}

var UndoFlip = DoFlip

func DoShuffle(ops []Operation, idx int, deckSize int, iterations int) int {
	for i := 0; i < iterations; i++ {
		for _, op := range ops {
			switch op.op {
			case OpFlip:
				idx = DoFlip(idx, deckSize)
			case OpDeal:
				idx = DoDeal(op.n, idx, deckSize)
			case OpCut:
				idx = DoCut(op.n, idx, deckSize)
			}
		}
	}
	return idx
}

func UndoShuffle(ops []Operation, idx int, deckSize int, iterations int) int {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		if i%1000000 == 0 {
			log.Print(int(float64(i) / time.Since(start).Seconds()))
		}
		for j := len(ops) - 1; j >= 0; j-- {
			switch ops[j].op {
			case OpFlip:
				idx = UndoFlip(idx, deckSize)
			case OpDeal:
				idx = UndoDeal(ops[j].n, idx, deckSize)
			case OpCut:
				idx = UndoCut(ops[j].n, idx, deckSize)
			}
		}
	}
	return idx
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

func Simplify(ops []Operation, deckSize int) []Operation {
	ds := big.NewInt(int64(deckSize))
	newOps := make([]Operation, len(ops))
	copy(newOps, ops)
	updated := true
loop:
	for updated {
		//log.Print(newOps)
		updated = false
		for i := 0; i < len(newOps); i++ {
			if newOps[i].op == OpCut && newOps[i].n < 0 {
				newOps[i].n += deckSize
			}
		}
		for i := 0; i < len(newOps)-1; i++ {
			if newOps[i].op == OpFlip && newOps[i+1].op == OpFlip {
				//log.Print("Remove both ", newOps[i], " and ", newOps[i+1])
				newOps = append(newOps[:i], newOps[i+2:]...)
				updated = true
				continue loop
			}
			if newOps[i].op == OpCut && newOps[i+1].op == OpFlip {
				//log.Print("Reverse ", newOps[i], " and ", newOps[i+1])
				newOps[i], newOps[i+1] = newOps[i+1], newOps[i]
				newOps[i+1].n = deckSize - newOps[i+1].n
				updated = true
				continue loop
			}
			if newOps[i].op == OpCut && newOps[i+1].op == OpCut {
				//log.Print("Combine ", newOps[i], " and ", newOps[i+1])
				i0 := big.NewInt(int64(newOps[i].n))
				i1 := big.NewInt(int64(newOps[i+1].n))
				i0.Add(i0, i1)
				i0.Mod(i0, ds)
				newOps[i].n = int(i0.Int64())
				newOps = append(newOps[:i+1], newOps[i+2:]...)
				updated = true
				continue loop
			}
			if newOps[i].op == OpDeal && newOps[i+1].op == OpFlip {
				//log.Print("Reverse ", newOps[i], " and ", newOps[i+1])
				newNewOps := make([]Operation, len(newOps)+1)
				copy(newNewOps, newOps[:i])
				newNewOps[i] = newOps[i+1]
				newNewOps[i+1] = newOps[i]
				newNewOps[i+2] = Operation{OpCut, -(newOps[i].n - 1)}
				copy(newNewOps[i+3:], newOps[i+2:])
				newOps = newNewOps
				updated = true
				continue loop
			}
			if newOps[i].op == OpDeal && newOps[i+1].op == OpCut {
				//log.Printf("reverse %v,%v", newOps[i], newOps[i+1])
				deal := newOps[i]
				cut := newOps[i+1]
				newOps[i] = Operation{OpCut, UndoDeal(deal.n, cut.n, deckSize)}
				newOps[i+1] = deal
				//log.Printf("%v,%v -> %v,%v", deal, cut, newOps[i], deal)
				updated = true
				continue loop
			}
			if newOps[i].op == OpDeal && newOps[i+1].op == OpDeal {
				//log.Printf("combine %v,%v", newOps[i], newOps[i+1])
				i0 := big.NewInt(int64(newOps[i].n))
				i1 := big.NewInt(int64(newOps[i+1].n))
				i0.Mul(i0, i1)
				i0.Mod(i0, ds)
				if !i0.IsInt64() {
					panic("result " + i0.String() + " overflows int64")
				}
				newOps[i].n = int(i0.Int64())
				newOps = append(newOps[:i+1], newOps[i+2:]...)
				updated = true
				continue loop
			}
		}
	}
	return newOps
}

func main() {
	ops, err := LoadShuffle("input")
	if err != nil {
		log.Fatal(err)
	}

	const iterations = 101741582076661
	const deckSize = 119315717514047

	base := Simplify(ops, deckSize)
	ops = make([]Operation, len(base))
	copy(ops, base)

	pow := map[int][]Operation{
		1: base,
	}
	i := 1
	for i*2 < iterations {
		i *= 2
		in := make([]Operation, 2*len(ops))
		copy(in, ops)
		copy(in[len(ops):], ops)
		ops = Simplify(in, deckSize) // double the iteration
		log.Printf("%d: %v -> %v", i, in, ops)
		pow[i] = make([]Operation, len(ops))
		copy(pow[i], ops)
	}

	powNum := i / 2
	for {
		if i == iterations {
			break
		}
		for i+powNum > iterations {
			powNum /= 2
		}
		in := make([]Operation, len(ops)+len(pow[powNum]))
		copy(in, ops)
		copy(in[len(ops):], pow[powNum])
		ops = Simplify(in, deckSize)
		i += powNum
		log.Printf("%d: %v -> %v", i, in, ops)
	}

	log.Print(UndoShuffle(ops, 2020, 119315717514047, 1))
}
