package main

import (
	"io/ioutil"
	"log"
	"strings"
)

type Point struct{ X, Y int }

type Map struct {
	Asteroids map[Point]bool
	Bounds    Point
}

func LoadMap(in string) Map {
	rows := strings.Split(strings.TrimSpace(in), "\n")
	ret := Map{
		Asteroids: map[Point]bool{},
		Bounds:    Point{len(rows[0]), len(rows)},
	}
	for x := 0; x < len(rows[0]); x++ {
		for y, row := range rows {
			if row[x] == '#' {
				ret.Asteroids[Point{x, y}] = true
			}
		}
	}

	return ret
}

func Slope(a, b Point) (dY, dX int) {
	dY = b.Y - a.Y // 6-2 == 4
	dX = b.X - a.X // 7-1 == 6

	absDY := dY
	if absDY < 0 {
		absDY *= -1
	}
	// absDY = 4

	absDX := dX
	if absDX < 0 {
		absDX *= -1
	}
	// absDX = 6

	if dY == 0 && dX == 0 {
		return 0, 0
	}
	if dY == 0 {
		return 0, dX / absDX
	}
	if dX == 0 {
		return dY / absDY, 0
	}

	// https://en.wikipedia.org/wiki/Euclidean_algorithm
	gcd := absDX
	k := absDY
	for k != 0 {
		gcd, k = k, gcd%k
	}

	dX /= gcd
	dY /= gcd

	return dY, dX
}

//func (m Map) OnLine(from Point, to Point) []Point {
//	dY, dX := Slope(from, to)
//	var ret []Point
//	cursor := Point{from.X + dX, from.Y + dY}
//	for cursor != to {
//		if m.Asteroids[cursor] {
//			ret = append(ret, cursor)
//		}
//		cursor.X += dX
//		cursor.Y += dY
//	}
//	return ret
//}

func (m Map) SightLines() map[Point]map[Point][]Point {
	ret := map[Point]map[Point][]Point{}
	for a := range m.Asteroids {
		ret[a] = map[Point][]Point{}
		for b := range m.Asteroids {
			if a == b {
				continue
			}
			s := Point{}
			s.Y, s.X = Slope(a, b)
			ret[a][s] = append(ret[a][s], b)
		}
	}
	return ret
}

func main() {
	input, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	m := LoadMap(string(input))
	lines := m.SightLines()
	var best Point
	var bestN int
	for point, lines := range lines {
		if len(lines) > bestN {
			best = point
			bestN = len(lines)
		}
	}
	log.Printf("Best point is %v wth %d lines", best, bestN)
}
