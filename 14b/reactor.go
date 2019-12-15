package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ReactionSet map[Material]Reaction

type Material string

type Reaction struct {
	Inputs      map[Material]int
	Output      Material
	OutputCount int
}

func (r Reaction) Perform(mats map[Material]int, target int) {
	mult := int(math.Ceil(float64(target) / float64(r.OutputCount)))
	for m, n := range r.Inputs {
		mats[m] -= n * mult
	}
	mats[r.Output] += r.OutputCount * mult
}

func ParseReactions(input string) (rs ReactionSet, err error) {
	rs = ReactionSet{}
	for _, reaction := range strings.Split(input, "\n") {
		log.Printf("parsing %q", reaction)
		react := Reaction{
			Inputs: map[Material]int{},
		}
		parts := strings.Split(reaction, "=>")
		if len(parts) != 2 {
			return nil, fmt.Errorf("reaction %q did not contain '=>'", reaction)
		}
		inputs := strings.Split(parts[0], ",")
		for _, input := range inputs {
			inputParts := strings.Split(strings.TrimSpace(input), " ")
			if len(inputParts) != 2 {
				return nil, fmt.Errorf("input %q didn't have two space-separated parts", input)
			}
			n, err := strconv.Atoi(inputParts[0])
			if err != nil {
				return nil, fmt.Errorf("input %q quantity was not an integer: %v", input, err)
			}
			react.Inputs[Material(inputParts[1])] = n
		}

		outputParts := strings.Split(strings.TrimSpace(parts[1]), " ")
		if len(outputParts) != 2 {
			return nil, fmt.Errorf("output %q didn't have two space-separated parts", parts[1])
		}
		outputN, err := strconv.Atoi(outputParts[0])
		if err != nil {
			return nil, fmt.Errorf("output quantity %q was not an integer", outputParts[0])
		}
		react.Output = Material(outputParts[1])
		react.OutputCount = outputN

		if _, ok := rs[react.Output]; ok {
			log.Fatalf("duplicate reactions producing the same output; %v and %v", rs[react.Output], react)
		}
		rs[react.Output] = react
	}
	return rs, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <filename>", os.Args[0])
	}
	reactionList, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("opening %q: %v", os.Args[1], err)
	}
	rs, err := ParseReactions(strings.TrimSpace(string(reactionList)))

	if err != nil {
		log.Fatalf("parsing reactions: %v", err)
	}

	min := 1
	max := int(1e12)
	for {
		amt := (min + max)
		if amt%2 == 0 {
			amt /= 2
		} else {
			amt = amt/2 + 1
		}
		log.Printf("%d, %d => %d", min, max, amt)
		have := map[Material]int{
			Material("FUEL"): -amt,
			Material("ORE"):  1e12,
		}

		// Scan through our materials, if we have a deficit of anything, perform
		// whatever reaction produces it. Do this until everything is positive or zero.
		done := false
		for !done {
			done = true
			for m, n := range have {
				if m == Material("ORE") {
					continue
				}
				if n < 0 {
					rs[m].Perform(have, -n)
					done = false
				}
			}
		}
		if have[Material("ORE")] < 0 {
			log.Print("high")
			max = amt-1
		}
		if have[Material("ORE")] >= 0 {
			log.Print("low")
			if amt == max {
				break
			}
			min = amt
		}
	}
	log.Print(min)
}
