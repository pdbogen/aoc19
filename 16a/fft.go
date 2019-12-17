package main

import (
	"bytes"
	"io/ioutil"
	"log"
)

var pattern = []int{0, 1, 0, -1}

// Each element in the new list is built by multiplying every value in the input
// list by a value in a repeating pattern and then adding up the results.
func FFT(in []int, phases int) []int {
	prePhase := make([]int, len(in))
	copy(prePhase, in)
	postPhase := make([]int, len(in))
	for i := 0; i < phases; i++ {
		for n := range in {
			postPhase[n] = FFTdigit(prePhase, n)
		}
		copy(prePhase, postPhase)
	}
	return postPhase
}

func FFTdigit(in []int, n int) (out int) {
	for j, k := range in {
		p := ((j + 1) / (n + 1)) % len(pattern)
		out += k * pattern[p]
	}
	if out < 0 {
		out *= -1
	}
	return out % 10
}

func main() {
	signalBytes, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("reading signal from \"input\": %v", err)
	}

	signalBytes = bytes.TrimSpace(signalBytes)
	signal := make([]int, len(signalBytes))
	for i, b := range signalBytes {
		signal[i] = int(b - '0')
	}

	log.Print(FFT(signal, 100))
}
