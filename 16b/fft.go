package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

// Each element in the new list is built by multiplying every value in the input
// list by a value in a repeating pattern and then adding up the results.
func FFT(in []int, phases int) []int {
	b := make([]int, len(in))
	for n := range in {
		b[n] = FFTdigit(in, n, phases)
	}
	return b
}

type cacheKey struct{ n, pass int }

var cache = map[cacheKey]int{}

func FFTdigit(in []int, n int, pass int) (out int) {
	if pass == 0 {
		return in[n]
	}
	if n == len(in)-1 {
		return in[n]
	}

	ck := cacheKey{n, pass}
	v, ok := cache[ck]
	if ok {
		return v
	}

	if n > len(in)/2 {
		out = FFTdigit(in, n, pass-1) + FFTdigit(in, n+1, pass)
	} else {
		idx := n // skip n items
		positive := true
	outer:
		for idx < len(in) {
			// handle n+1 numbers with mult p
			// skip i numbers
			// handle n+1 numbers with mult -p
			// skip i numbers
			for span := 0; span < n+1; span++ {
				if idx+span >= len(in) {
					break outer
				}
				if positive {
					out += FFTdigit(in, idx+span, pass-1)
				} else {
					out -= FFTdigit(in, idx+span, pass-1)
				}
			}
			idx += n + 1
			positive = !positive
			idx += n + 1 // skip n+1
		}
	}

	if out < 0 {
		out *= -1
	}
	out = out % 10

	cache[ck] = out
	return out
}

func Repeat(in []int, n int) []int {
	out := make([]int, len(in)*n)
	for i := 0; i < n; i++ {
		for j, k := range in {
			out[j+len(in)*i] = k
		}
	}
	return out
}

func Decode(signal []int) []int {
	offset := 0
	for i := 0; i < 7; i++ {
		offset = 10*offset + signal[i]
	}
	fmt.Println(offset)

	repsignal := Repeat(signal, 10000)
	result := make([]int, 8)
	for i := offset; i < offset+8; i++ {
		result[i-offset] = FFTdigit(repsignal, i, 100)
	}
	return result
}

func main() {
	//fmt.Print(FFTdigit([]int{1, 2, 3, 4, 5, 6, 7, 8}, 1, 0))
	//fmt.Print(FFTdigit([]int{1, 2, 3, 4, 5, 6, 7, 8}, 1, 1))
	signalBytes, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("reading signal from \"input\": %v", err)
	}

	signalBytes = bytes.TrimSpace(signalBytes)
	signal := make([]int, len(signalBytes))
	for i, b := range signalBytes {
		signal[i] = int(b - '0')
	}

	log.Print(Decode(signal))
	//log.Print(Decode([]int{0, 3, 0, 3, 6, 7, 3, 2, 5, 7, 7, 2, 1, 2, 9, 4, 4, 0, 6, 3, 4, 9, 1, 5, 6, 5, 4, 7, 4, 6, 6, 4}))
	//log.Print(Decode([]int{0, 2, 9, 3, 5, 1, 0, 9, 6, 9, 9, 9, 4, 0, 8, 0, 7, 4, 0, 7, 5, 8, 5, 4, 4, 7, 0, 3, 4, 3, 2, 3}))
}
