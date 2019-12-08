package main

import (
	"github.com/pdbogen/aoc19/sif"
	"io/ioutil"
	"log"
	"math"
)

func countDigits(layers [][][]int, layer int, digit int) int {
	count := 0
	for x := range layers {
		for y := range layers[0] {
			if layers[x][y][layer] == digit {
				count++
			}
		}
	}
	return count
}

func main() {
	data, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("could not parse input: %v", err)
	}

	log.Printf("%d total layers", len(data)/25/6)
	layers := sif.ToLayers(string(data), 25, 6)
	minZeroes := math.MaxInt32
	minZeroLayer := 0
	for l := range layers[0][0] {
		zeroes := countDigits(layers, l, 0)
		log.Printf("%d zeroes on layer %d", zeroes, l)
		if zeroes < minZeroes {
			minZeroes = zeroes
			minZeroLayer = l
		}
	}
	log.Printf("%d zeroes on layer %d", minZeroes, minZeroLayer)
	log.Printf("layer %d checksum %d", minZeroLayer,
		countDigits(layers, minZeroLayer, 1)*countDigits(layers, minZeroLayer, 2),
	)
}
