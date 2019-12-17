package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"log"
)

type Point struct{ X, Y int }

func Print(scaffold map[Point]byte) {
	x, y := 0, 0
	for {
		c, ok := scaffold[Point{x, y}]
		if !ok {
			if x == 0 {
				break
			}
			x = 0
			y++
			fmt.Printf("\n")
			continue
		}
		fmt.Printf("%c", c)
		x++
	}
}

func Intersections(scaffold map[Point]byte) (ret []Point) {
scaffold:
	for p, c := range scaffold {
		if c != 35 {
			continue
		}
		for _, cardinal := range []Point{
			{p.X - 1, p.Y},
			{p.X + 1, p.Y},
			{p.X, p.Y - 1},
			{p.X, p.Y + 1},
		} {
			if scaffold[cardinal] != 35 {
				continue scaffold
			}
		}
		ret = append(ret, p)
	}
	return ret
}

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	inCh, outCh := make(chan int), make(chan int)
	go intcode.Execute(prog, inCh, outCh)
	scaffold := map[Point]byte{}
	x := 0
	y := 0
	for c := range outCh {
		if c == 10 {
			x = 0
			y++
			continue
		}
		scaffold[Point{x, y}] = byte(c)
		x++
	}
	Print(scaffold)
	//log.Print(scaffold)
	alignSum := 0
	for _, isect := range Intersections(scaffold) {
		alignSum += isect.X * isect.Y
	}
	log.Printf("Calibratin Value: %d", alignSum)
}
