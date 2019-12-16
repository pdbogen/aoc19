package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"github.com/pdbogen/aoc19/term"
	"log"
	"time"
)

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

type Point struct{ X, Y int }

type World struct {
	in      chan<- int
	out     <-chan int
	Drone   Point
	Map     map[Point]rune
	lastMin Point
	lastMax Point
}

func (w *World) Print() {
	var min, max Point
	for pt := range w.Map {
		if pt.X < min.X {
			min.X = pt.X
		}
		if pt.X > max.X {
			max.X = pt.X
		}
		if pt.Y < min.Y {
			min.Y = pt.Y
		}
		if pt.Y > max.Y {
			max.Y = pt.Y
		}
	}
	if w.lastMin != min || w.lastMax != max {
		term.Clear()
	}
	w.lastMin = min
	w.lastMax = max
	term.MoveCursor(1, 1)
	for y := max.Y; y >= min.Y; y-- {
		var line []rune
		for x := min.X; x <= max.X; x++ {
			pt := Point{x, y}
			c, ok := w.Map[Point{x, y}]
			if pt == w.Drone && c != '@' {
				line = append(line, ([]rune(term.Scolor(85, 255, 255)))...)
				line = append(line, 'D')
				line = append(line, ([]rune(term.ScolorReset()))...)
				continue
			}
			if ok {
				if c == '@' {
					line = append(line, ([]rune(term.Scolor(85, 255, 85)))...)
				} else if c == '#' {
					line = append(line, ([]rune(term.Scolor(170, 85, 0)))...)
				} else if c == '+' {
					line = append(line, ([]rune(term.Scolor(255, 85, 255)))...)
				}
				line = append(line, c)
				line = append(line, ([]rune(term.ScolorReset()))...)
			} else {
				line = append(line, ' ')
			}
		}
		fmt.Println(string(line))
	}
}

func (w *World) Coord(dir int) Point {
	switch (dir) {
	case North:
		return Point{w.Drone.X, w.Drone.Y + 1}
	case South:
		return Point{w.Drone.X, w.Drone.Y - 1}
	case East:
		return Point{w.Drone.X + 1, w.Drone.Y}
	case West:
		fallthrough
	default:
		return Point{w.Drone.X - 1, w.Drone.Y}
	}
}

func (w *World) Check() {
	for dir, back := range map[int]int{
		North: South,
		East:  West,
		South: North,
		West:  East,
	} {
		if _, ok := w.Map[w.Coord(dir)]; ok {
			continue
		}
		w.in <- dir
		result := <-w.out
		if result == 0 {
			w.Map[w.Coord(dir)] = '#'
		} else if result == 1 {
			w.Map[w.Coord(dir)] = '+'
			w.in <- back
			<-w.out
		} else if result == 2 {
			w.Map[w.Coord(dir)] = '@'
			w.in <- back
			<-w.out
		}
	}
}

func (w *World) Go(dir int) bool {
	time.Sleep(5 * time.Millisecond)
	w.in <- dir
	result := <-w.out
	newPt := w.Coord(dir)
	if result == 0 {
		w.Map[newPt] = '#'
	} else if result == 1 {
		w.Map[newPt] = '.'
		w.Drone = newPt
		w.Check()
	} else if result == 2 {
		w.Map[newPt] = '@'
		w.Drone = newPt
		w.Check()
	}
	return !(result == 0)
}

func Reverse(dir int) int {
	switch (dir) {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}
	return -1
}

// Explore traverses the map in oen of two modes:
// 1. if target is zero, Explore finds the farthest possible square from the
// current location, traversing the entire map in the process. It returns the
// distance to this location.
// 2. if the target is any other rune, Explore explores the map until finding
// this square, and then stops, returning 1 if the square is found.
func (w *World) Explore(dir int, steps int, target rune) (maxSteps int) {
	maxSteps = steps
	if w.Go(dir) {
		w.Print()
		for _, nextDir := range []int{North, East, South, West} {
			if nextDir == Reverse(dir) {
				continue
			}
			nextSteps := w.Explore(nextDir, steps+1, target)
			if target != 0 && w.Map[w.Drone] == target {
				return 1
			}
			if target == 0 && nextSteps > maxSteps {
				maxSteps = nextSteps
			}
		}
		w.Go(Reverse(dir))
		w.Print()
	}
	return maxSteps
}

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatalf("loading program: %v", err)
	}

	inCh := make(chan int)
	outCh := make(chan int)
	go func() {
		if _, err := intcode.Execute(prog, inCh, outCh); err != nil {
			log.Fatal(err)
		}
	}()

	world := World{
		in:  inCh,
		out: outCh,
		Map: map[Point]rune{},
	}

	term.HideCursor()
	for _, dir := range []int{North, East, South, West} {
		if world.Explore(dir, 0, '@') > 0 {
			break
		}
	}

	var maxSteps int
	for _, dir := range []int{North, East, South, West} {
		nextSteps := world.Explore(dir, 0, rune(0))
		if nextSteps > maxSteps {
			maxSteps = nextSteps
		}
	}
	term.ShowCursor()
	log.Printf("max steps: %d", maxSteps)
}
