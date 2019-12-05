package main

import (
	"fmt"
	"strconv"
)

type Rule func(in int, inS string) bool

func NumbersIncreaseToTheLeft(_ int, inS string) bool {
	last := '0'
	for _, c := range inS {
		if c < last {
			return false
		}
		last = c
	}
	return true
}

func AtLeastOneExactDouble(_ int, inS string) bool {
	group := ""
	for _, c := range inS {
		if len(group) == 0 {
			group = string(c)
			continue
		}
		if []rune(group)[0] == c {
			group = group + string(c)
			continue
		}
		if len(group) == 2 {
			return true
		}
		group = string(c)
	}
	return len(group) == 2
}

func main() {
	rules := []Rule{
		NumbersIncreaseToTheLeft,
		AtLeastOneExactDouble,
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
