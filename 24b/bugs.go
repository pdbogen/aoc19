package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/term"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Point struct {
	X, Y, Level int
}

func (p Point) Adjacent(width int) []Point {
	var ret []Point
	if p.X == 0 {
		ret = append(ret,
			Point{width/2 - 1, width / 2, p.Level - 1},
			Point{p.X + 1, p.Y, p.Level})
	} else if p.X == (width - 1) {
		ret = append(ret,
			Point{p.X - 1, p.Y, p.Level},
			Point{width/2 + 1, width / 2, p.Level - 1})
	} else if p.Y == width/2 && p.X == width/2-1 {
		ret = append(ret, Point{p.X - 1, p.Y, p.Level})
		for y := 0; y < width; y++ {
			ret = append(ret, Point{0, y, p.Level + 1})
		}
	} else if p.Y == 2 && p.X == width/2+1 {
		for y := 0; y < width; y++ {
			ret = append(ret, Point{width - 1, y, p.Level + 1})
		}
		ret = append(ret, Point{p.X + 1, p.Y, p.Level})
	} else {
		ret = append(ret,
			Point{p.X - 1, p.Y, p.Level},
			Point{p.X + 1, p.Y, p.Level},
		)
	}
	if p.Y == 0 {
		ret = append(ret,
			Point{width / 2, width/2 - 1, p.Level - 1},
			Point{p.X, p.Y + 1, p.Level})
	} else if p.Y == (width - 1) {
		ret = append(ret,
			Point{width / 2, width/2 + 1, p.Level - 1},
			Point{p.X, p.Y - 1, p.Level})
	} else if p.X == width/2 && p.Y == width/2-1 {
		ret = append(ret, Point{p.X, p.Y - 1, p.Level})
		for x := 0; x < width; x++ {
			ret = append(ret, Point{x, 0, p.Level + 1})
		}
	} else if p.X == width/2 && p.Y == width/2+1 {
		for x := 0; x < width; x++ {
			ret = append(ret, Point{x, width - 1, p.Level + 1})
		}
		ret = append(ret, Point{p.X, p.Y + 1, p.Level})
	} else {
		ret = append(ret,
			Point{p.X, p.Y - 1, p.Level},
			Point{p.X, p.Y + 1, p.Level})
	}
	return ret
}

type Space struct {
	Bug          bool
	AdjacentBugs int
}

type Grid struct {
	Width    int
	Grid     map[Point]*Space
	Min, Max Point
}

func (g *Grid) Count() {
	var points []Point
	for pt, space := range g.Grid {
		space.AdjacentBugs = 0
		if space.Bug {
			points = append(points, pt)
		}
	}
	for _, pt := range points {
		for _, adj := range pt.Adjacent(g.Width) {
			if _, ok := g.Grid[adj]; !ok {
				g.Grid[adj] = &Space{}
			}
			g.Grid[adj].AdjacentBugs++
		}
	}
	for pt := range g.Grid {
		if pt.X < g.Min.X {
			g.Min = Point{pt.X, g.Min.Y, g.Min.Level}
		}
		if pt.X > g.Max.X {
			g.Max = Point{pt.X, g.Max.Y, g.Max.Level}
		}
		if pt.Y < g.Min.Y {
			g.Min = Point{g.Min.X, pt.Y, g.Min.Level}
		}
		if pt.Y > g.Max.Y {
			g.Max = Point{g.Max.X, pt.Y, g.Max.Level}
		}
		if pt.Level < g.Min.Level {
			g.Min = Point{g.Min.X, g.Min.Y, pt.Level}
		}
		if pt.Level > g.Max.Level {
			g.Max = Point{g.Max.X, g.Max.Y, pt.Level}
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
}

func (g Grid) String() string {
	ret := ""
	for level := g.Min.Level; level <= g.Max.Level; level++ {
		ret += fmt.Sprintf(fmt.Sprintf("%%%dd ", g.Width), level)
	}
	ret += "\n"
	for y := g.Min.Y; y <= g.Max.Y; y++ {
		for level := g.Min.Level; level <= g.Max.Level; level++ {
			for x := g.Min.X; x <= g.Max.X; x++ {
				space, ok := g.Grid[Point{x, y, level}]
				if ok && space.Bug {
					ret += "#"
				} else {
					ret += "."
				}
			}
			ret += " "
		}
		ret += "\n"
	}
	bugs := 0
	for _, s := range g.Grid {
		if s.Bug {
			bugs++
		}
	}
	ret += fmt.Sprintf("%d bugs on grid\n", bugs)
	return ret
}

func Load(in io.Reader) *Grid {
	ret := &Grid{
		Grid:  map[Point]*Space{},
		Width: 5,
	}
	input, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatalf("err loading input: %v", err)
	}

	for y, line := range strings.Split(string(input), "\n") {
		for x, c := range line {
			ret.Grid[Point{x, y, 0}] = &Space{
				Bug: c == '#',
			}
		}
	}
	return ret
}

//func (g Grid) Rating() int {
//	ret := 0
//	for pt, space := range g.Grid {
//		ptv := int(math.Pow(2, float64(pt.Y*(g.Bounds.X+1)+pt.X)))
//		if space.Bug {
//			ret += ptv
//		}
//	}
//	return ret
//}

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatalf("loading input: %v", err)
	}
	grid := Load(input)
//	grid := Load(bytes.NewBufferString(`....#
//#..#.
//#..##
//..#..
//#....`))
	for i := 0; i < 200; i++ {
		grid.Count()
		grid.Update()
		grid.Log()
	}
	//term.Clear()
	//grid.Print()
}
