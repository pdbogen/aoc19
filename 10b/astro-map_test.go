package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSlope(t *testing.T) {
	type args struct {
		a Point
		b Point
	}
	tests := []struct {
		a, b   Point
		wantDY int
		wantDX int
	}{
		{Point{0, 0}, Point{1, 1}, 1, 1},
		{Point{0, 0}, Point{2, 6}, 3, 1},
		{Point{0, 0}, Point{6, 2}, 1, 3},
		{Point{0, 0}, Point{-2, 6}, 3, -1},
		{Point{0, 0}, Point{2, -6}, -3, 1},
		{Point{0, 0}, Point{-2, -6}, -3, -1},
		{Point{0, 0}, Point{-6, 2}, 1, -3},
		{Point{0, 0}, Point{6, -2}, -1, 3},
		{Point{0, 0}, Point{-6, -2}, -1, -3},
		{Point{1, 2}, Point{7, 6}, 2, 3},
		{Point{3, 4}, Point{1, 0}, -2, -1},
	}
	for _, tt := range tests {
		gotDY, gotDX := Slope(tt.a, tt.b)
		if gotDY != tt.wantDY {
			t.Errorf("Slope(%v,%v) gotDY = %v, want %v", tt.a, tt.b, gotDY, tt.wantDY)
		}
		if gotDX != tt.wantDX {
			t.Errorf("Slope(%v,%v) gotDX = %v, want %v", tt.a, tt.b, gotDX, tt.wantDX)
		}
	}
}

type TestMap struct {
	Name      string
	Map       string
	Best      Point
	BestCount int
}

var Example2 = []string{
	"#.#...#.#.",
	".###....#.",
	".#....#...",
	"##.#.#.#.#",
	"....#.#.#.",
	".##..###.#",
	"..#...##..",
	"..##....##",
	"......#...",
	".####.###.",
}

const Example3 = `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`

const Example4 = `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`

var Tests = []TestMap{
	{"example 0",
		`.#..#
.....
#####
....#
...##`, Point{3, 4}, 8},
	{"example 1",
		`......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`, Point{5, 8}, 33},
	{"example 2a", strings.Join(Example2[0:1], "\n"), Point{2, 0}, 2},
	{"example 2b", strings.Join(Example2[0:2], "\n"), Point{2, 1}, 6},
	{"example 2c", strings.Join(Example2[0:3], "\n"), Point{1, 2}, 9},
	{"example 2d", strings.Join(Example2[0:4], "\n"), Point{1, 2}, 15},
	{"example 2e", strings.Join(Example2[0:5], "\n"), Point{1, 2}, 18},
	{"example 2", strings.Join(Example2, "\n"), Point{1, 2}, 35},
	{"example 3", Example3, Point{6, 3}, 41},
	{"example 4", Example4, Point{11, 13}, 210},
}

//func TestMap_OnLine(t *testing.T) {
//	tests := []struct {
//		name string
//		m    Map
//		from Point
//		to   Point
//		want []Point
//	}{
//		{
//			"two roid",
//			Map{
//				Asteroids: map[Point]bool{
//					Point{0, 0,}: true,
//					Point{1, 1}:  true,
//				},
//			},
//			Point{0, 0},
//			Point{1, 1},
//			nil,
//		},
//		{
//			"three roid, no isect, a - b",
//			Map{
//				Asteroids: map[Point]bool{
//					Point{0, 0,}: true,
//					Point{1, 1}:  true,
//					Point{1, 2}:  true,
//				},
//			},
//			Point{0, 0},
//			Point{1, 1},
//			nil,
//		},
//		{
//			"three roid, no isect, a - c",
//			Map{
//				Asteroids: map[Point]bool{
//					Point{0, 0}: true,
//					Point{1, 1}: true,
//					Point{1, 2}: true,
//				},
//			},
//			Point{0, 0},
//			Point{1, 2},
//			nil,
//		},
//		{
//			"three roid, isect, a - b",
//			Map{
//				Asteroids: map[Point]bool{
//					Point{0, 0,}: true,
//					Point{1, 1}:  true,
//					Point{2, 2}:  true,
//				},
//			},
//			Point{0, 0},
//			Point{1, 1},
//			nil,
//		},
//		{
//			"three roid, isect, a - c",
//			Map{
//				Asteroids: map[Point]bool{
//					Point{0, 0}: true,
//					Point{1, 1}: true,
//					Point{2, 2}: true,
//				},
//			},
//			Point{0, 0},
//			Point{2, 2},
//			[]Point{{1, 1}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := tt.m.OnLine(tt.from, tt.to)
//			assert.Equal(t, len(tt.want), len(got))
//			if tt.want == nil {
//				assert.Empty(t, got)
//				return
//			}
//			assert.Subset(t, got, tt.want)
//		})
//	}
//}

func TestMap_SightLines(t *testing.T) {
	for _, test := range Tests {
		t.Run(test.Name, func(t *testing.T) {
			m := LoadMap(test.Map)
			bases := m.SightLines()
			if !assert.Equal(t, test.BestCount, len(bases[test.Best])) {
				for pt, isects := range bases[test.Best] {
					t.Logf("%v -> %v", pt, isects)
				}
			}
		})
	}
}
