package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct{ X, Y int }
type Move struct {
	Direction Direction
	Distance  int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var Directions = map[rune]Direction{
	'U': Up,
	'D': Down,
	'L': Left,
	'R': Right,
}

func ParseLine(line string) (ret []Move, err error) {
	for _, move := range strings.Split(line, ",") {
		move = strings.TrimSpace(move)
		if len(move) < 2 {
			return nil, fmt.Errorf("move %q too short", move)
		}
		dirRune := ([]rune(move))[0]
		dir, ok := Directions[dirRune]
		if !ok {
			return nil, fmt.Errorf("unrecognized direction %c", dirRune)
		}

		dist, err := strconv.Atoi(move[1:])
		if err != nil {
			return nil, fmt.Errorf("could not parse move distance %q: %v", move[1:], err)
		}

		ret = append(ret, Move{dir, dist,})
	}

	return ret, nil
}

func MustToGrid(lines []string) (Grid) {
	g, e := ToGrid(lines)
	if e != nil {
		log.Fatalf("MustToGrid failed to convert %v: %v", lines, e)
	}
	return g
}

func ToGrid(lines []string) (Grid, error) {
	grid := Grid{}
	for lineNo, line := range lines {
		var pen Point
		steps := 0
		moves, err := ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parsing line %q: %v", line, err)
		}
		for _, move := range moves {
			xDir, yDir := 0, 0
			xFin, yFin := pen.X, pen.Y
			switch move.Direction {
			case Up:
				yDir = 1
				yFin += move.Distance
			case Down:
				yDir = -1
				yFin -= move.Distance
			case Right:
				xDir = 1
				xFin += move.Distance
			case Left:
				xDir = -1
				xFin -= move.Distance
			}
			for {
				pen.X += xDir
				pen.Y += yDir
				steps++
				if _, ok := grid[pen]; !ok {
					grid[pen] = map[int]int{}
				}
				grid[pen][lineNo] = steps
				if pen.X == xFin && pen.Y == yFin {
					break
				}
			}
		}
	}
	return grid, nil
}

func FindCrossings(grid Grid) (map[Point]int, error) {
	points := map[Point]int{}

	for pt, wires := range grid {
		if len(wires) >= 2 {
			points[pt] = wires[0] + wires[1]
		}
	}

	return points, nil
}

func ShortestDelayCrossing(grid Grid) (Point, int, error) {
	crossings, err := FindCrossings(grid)
	if err != nil {
		return Point{}, 0, fmt.Errorf("finding crossings: %v", err)
	}

	if len(crossings) == 0 {
		return Point{}, 0, errors.New("no intersections!")
	}

	crossingList := make([]Point, 0, len(crossings))
	for pt := range crossings {
		crossingList = append(crossingList, pt)
	}

	sort.Slice(crossingList, func(i, j int) bool {
		return crossings[crossingList[i]] < crossings[crossingList[j]]
	})

	return crossingList[0], crossings[crossingList[0]], nil
}

func ClosestCrossingDistance(grid Grid) (int, error) {
	crossings, err := FindCrossings(grid)
	if err != nil {
		return 0, fmt.Errorf("finding crossings: %v", err)
	}

	if len(crossings) == 0 {
		return 0, errors.New("no intersections!")
	}

	crossingList := make([]Point, 0, len(crossings))
	for pt := range crossings {
		crossingList = append(crossingList, pt)
	}

	sort.Slice(crossingList, func(i, j int) bool {
		return dist2d(Point{}, crossingList[i]) < dist2d(Point{}, crossingList[j])
	})

	return dist2d(Point{}, crossingList[0]), nil
}

// manhattan distance
func dist2d(a, b Point) int {
	if a.X < b.X {
		a.X, b.X = b.X, a.X
	}
	if a.Y < b.Y {
		a.Y, b.Y = b.Y, a.Y
	}
	return a.X - b.X + a.Y - b.Y
}

func main() {
	showgrid := flag.Bool("grid", false, "if true, output grid.png")
	flag.Parse()

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(input)), "\n")

	if len(lines) != 2 {
		log.Fatalf("expected two wires to trace, but got %d", len(lines))
	}

	grid, err := ToGrid(lines)
	if err != nil {
		log.Fatalf("drawing lines onto grid: %v", err)
	}

	dist, err := ClosestCrossingDistance(grid)
	if err != nil {
		log.Fatalf("finding closest intersection: %v", err)
	}

	if *showgrid {
		log.Print("Saving image...")
		if err := grid.Save("grid.png"); err != nil {
			log.Fatal(err)
		}
		log.Print("Done")
	}
	log.Printf("Closest intersection distance: %d", dist)
	shortestPt, delay, err := ShortestDelayCrossing(grid)
	if err != nil {
		log.Fatalf("finding shortest delay: %v", err)
	}
	log.Printf("Shortest delay distance: %d at %v", delay, shortestPt)
}
