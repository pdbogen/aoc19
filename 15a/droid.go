package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"log"
)

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

type Point struct{ X, Y int }

type World struct {
	in    chan<- int
	out   <-chan int
	Drone Point
	Map   map[Point]rune
}

func (w World) Print() {
	minx, maxx, miny, maxy := w.Drone.X, w.Drone.X, w.Drone.Y, w.Drone.Y
	for pt := range w.Map {
		if pt.X < minx {
			minx = pt.X
		}
		if pt.X > maxx {
			maxx = pt.X
		}
		if pt.Y < miny {
			miny = pt.Y
		}
		if pt.Y > maxy {
			maxy = pt.Y
		}
	}
	log.Print("-----")
	for y := maxy; y >= miny; y-- {
		var line []rune
		for x := minx; x <= maxx; x++ {
			pt := Point{x, y}
			c, ok := w.Map[Point{x, y}]
			if pt == w.Drone && c != '@' {
				line = append(line, 'D')
				continue
			}
			if ok {
				line = append(line, c)
			} else {
				line = append(line, ' ')
			}
		}
		log.Print(string(line))
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
	w.in <- dir
	result := <-w.out
	newPt := w.Coord(dir)
	if result == 0 {
		w.Map[newPt] = '#'
	} else if result == 1 {
		if w.Map[w.Coord(dir)] == '.' {
			w.Map[w.Drone] = '+'
		}
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

func run() error {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		return fmt.Errorf("loading program: %v", err)
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
		Map: map[Point]rune{
			Point{0,0}: '.',
		},
	}

	dir := North
	for {
		// try to go left
		// else straight
		// else 180 & try above two again

		if !world.Go(Left(dir)) {
			dir = Right(dir)
		} else {
			dir = Left(dir)
		}

		if world.Map[world.Drone] == '@' {
			log.Printf("found at %v", world.Drone)
			break
		}
		world.Print()
	}
	n := 0
	for _, c := range world.Map {
		if c == '.' {
			n++
		}
	}
	log.Printf("%d steps", n)
	return nil
}

func Left(in int) int {
	switch in {
	case North:
		return West
	case West:
		return South
	case South:
		return East
	case East:
		return North
	}
	return -1
}

func Right(in int) int {
	switch (in) {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	}
	return -1
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
