package main

import (
	"fmt"
	"strconv"
)

type Rule func(in int, inS string) bool

func main() {
	rules := []Rule{
		func(in int, inS string) bool {
			last := '0'
			for _, c := range inS {
				if c < last {
					return false
				}
				last = c
			}
			return true
		},
		func(in int, inS string) bool {
			for i := 0; i < len(inS)-1; i++ {
				if inS[i] == inS[i+1] {
					return true
				}
			}
			return false
		},
	}

	potentials := []int{}
guessing:
	for i := 264793; i <= 803935; i++ {
		iS := strconv.Itoa(i)
		for _, rule := range rules {
			if !rule(i, iS) {
				continue guessing
			}
		}
		potentials = append(potentials, i)
		fmt.Println(i)
	}
	fmt.Printf("%d potential passwords\n", len(potentials))
}
