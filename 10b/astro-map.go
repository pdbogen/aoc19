package main

import (
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
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

	targets := lines[best]
	for _, tgt := range targets {
		sort.Slice(tgt, func(i, j int) bool {
			dxi := best.X - tgt[i].X
			dxi *= dxi
			dyi := best.Y - tgt[i].Y
			dyi *= dyi
			dxj := best.X - tgt[j].X
			dxj *= dxj
			dyj := best.Y - tgt[j].Y
			dyj *= dyj
			return (dxi + dyi) < (dxj + dyj)
		})
	}

	anim := &gif.GIF{}

	n := 0
	for {
		var points []Point
		for pt, atPt := range targets {
			if len(atPt) > 0 {
				points = append(points, pt)
			}
		}
		if len(points) == 0 {
			break
		}
		sort.Slice(points, func(i, j int) bool {
			ptI := points[i]
			ptJ := points[j]

			angleI := 90+math.Atan2(float64(ptI.Y), float64(ptI.X)) * 180 / math.Pi
			if angleI > 360 {
				angleI -= 360
			}
			if angleI < 0 {
				angleI += 360
			}

			angleJ := 90+math.Atan2(float64(ptJ.Y), float64(ptJ.X)) * 180 / math.Pi
			if angleJ > 360 {
				angleJ -= 360
			}
			if angleJ < 0 {
				angleJ += 360
			}

			// 90…0 -> 0…90
			// -180…0 -> 91…270
			// 90…180 -> 360…270
			return angleI < angleJ
		})
		for _, pt := range points {
			n++
			log.Printf("%d: %v at slope %v theta %f", n, targets[pt][0], pt,
				math.Atan2(float64(pt.Y), float64(pt.X))*180/math.Pi)
			tgt := targets[pt][0]
			targets[pt] = targets[pt][1:]
			anim.Image = append(anim.Image, Draw(best, m.Bounds, targets, tgt))
			anim.Delay = append(anim.Delay, 10)
			anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		}
	}

	f, _ := os.OpenFile("output.gif", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err := gif.EncodeAll(f, anim); err != nil {
		log.Fatal(err)
	}
	f.Close()
}
