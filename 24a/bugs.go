package main

import (
	"github.com/pdbogen/aoc19/common"
	"github.com/pdbogen/aoc19/term"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

type Space struct {
	Bug          bool
	AdjacentBugs int
}

type Grid struct {
	Grid   map[common.Point]*Space
	Bounds common.Point
}

func (g Grid) Count() {
	for pt, space := range g.Grid {
		space.AdjacentBugs = 0
		for _, apt := range []common.Point{pt.North(), pt.East(), pt.South(), pt.West()} {
			if a, ok := g.Grid[apt]; ok && a.Bug {
				space.AdjacentBugs++
			}
		}
	}
}

func (g Grid) Update() {
	for _, space := range g.Grid {
		if space.Bug && space.AdjacentBugs != 1 {
			space.Bug = false
		} else if space.AdjacentBugs == 1 || space.AdjacentBugs == 2 {
			space.Bug = true
		}
	}
}

func (g Grid) Log() {
	for _, line := range strings.Split(strings.TrimSpace(g.String()), "\n") {
		log.Print(line)
	}
}

func (g Grid) Print() {
	term.MoveCursor(1, 1)
	print(g.String())
	println(g.Rating())
}

func (g Grid) String() string {
	ret := ""
	for y := 0; y <= g.Bounds.Y; y++ {
		for x := 0; x <= g.Bounds.X; x++ {
			space, ok := g.Grid[common.Point{x, y}]
			if ok {
				if space.Bug {
					ret += "#"
				} else {
					ret += "."
				}
			} else {
				ret += " "
			}
		}
		ret += "\n"
	}
	return ret
}

func Load(in io.Reader) *Grid {
	ret := &Grid{
		Grid:   map[common.Point]*Space{},
		Bounds: common.Point{},
	}
	input, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatalf("err loading input: %v", err)
	}

	for y, line := range strings.Split(string(input), "\n") {
		if ret.Bounds.Y < y {
			ret.Bounds.Y = y
		}
		for x, c := range line {
			if ret.Bounds.X < x {
				ret.Bounds.X = x
			}
			ret.Grid[common.Point{x, y}] = &Space{
				Bug: c == '#',
			}
		}
		y++
	}
	return ret
}

func (g Grid) Rating() int {
	ret := 0
	for pt, space := range g.Grid {
		ptv := int(math.Pow(2, float64(pt.Y*(g.Bounds.X+1)+pt.X)))
		if space.Bug {
			ret += ptv
		}
	}
	return ret
}

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatalf("loading input: %v", err)
	}
	grid := Load(input)
	seen := map[string]bool{}
	term.Clear()
	for {
		grid.Print()
		s := grid.String()
		if seen[s] {
			log.Print("seen before!")
			os.Exit(0)
		}
		seen[s] = true
		grid.Count()
		grid.Update()
	}
}
