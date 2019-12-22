package main

import (
	"bytes"
	"fmt"
	"github.com/pdbogen/aoc19/common"
	"github.com/pdbogen/aoc19/term"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
)

type Map struct {
	Map       map[common.Point]*common.Location
	Wall      map[common.Point]bool
	Bounds    common.Point
	Portals   map[string][]*common.Location
	Start     *common.Location
	Goal      *common.Location
	pathCache *sync.Map
}

func LoadMap(reader io.Reader) (*Map, error) {
	ret := &Map{
		Map:       map[common.Point]*common.Location{},
		Portals:   map[string][]*common.Location{},
		pathCache: &sync.Map{},
		Wall:      map[common.Point]bool{},
	}

	grid, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	points := map[common.Point]byte{}

	// Locate starting positions, so that we can also identify sections.
	rows := bytes.Split(grid, []byte{'\n'})
	for y, row := range rows {
		for x, c := range row {
			points[common.Point{x, y}] = c
			if x > ret.Bounds.X {
				ret.Bounds.X = x
			}
			if y > ret.Bounds.Y {
				ret.Bounds.Y = y
			}
		}
	}
	for pt, c := range points {
		if c == '#' {
			ret.Wall[pt] = true
		}
		if c != '.' {
			continue
		}
		l := &common.Location{
			Point:       pt,
			Connections: map[common.Direction]*common.Location{},
		}
		ret.Map[pt] = l

		n, n2 := points[pt.North().North()], points[pt.North()]
		e, e2 := points[pt.East()], points[pt.East().East()]
		w, w2 := points[pt.West().West()], points[pt.West()]
		s, s2 := points[pt.South()], points[pt.South().South()]
		if n2 >= 'A' && n2 <= 'Z' {
			l.Trait = string([]byte{n, n2})
		}
		if e >= 'A' && e <= 'Z' {
			l.Trait = string([]byte{e, e2})
		}
		if w2 >= 'A' && w2 <= 'Z' {
			l.Trait = string([]byte{w, w2})
		}
		if s >= 'A' && s <= 'Z' {
			l.Trait = string([]byte{s, s2})
		}

		if l.Trait == "AA" {
			ret.Start = l
		} else if l.Trait == "ZZ" {
			ret.Goal = l
		} else if l.Trait != "" {
			ret.Portals[l.Trait] = append(ret.Portals[l.Trait], l)
			if len(ret.Portals[l.Trait]) == 2 {
				ret.Portals[l.Trait][0].Connections[common.None] = ret.Portals[l.Trait][1]
				ret.Portals[l.Trait][1].Connections[common.None] = ret.Portals[l.Trait][0]
			}
		}

		for d, p := range map[common.Direction]common.Point{
			common.North: pt.North(),
			common.East:  pt.East(),
			common.South: pt.South(),
			common.West:  pt.West(),
		} {
			if cl, ok := ret.Map[p]; ok {
				l.Connections[d] = cl
				cl.Connections[d.Opposite()] = l
			}
		}
	}

	return ret, nil
}

func (m *Map) Path(from, to common.Point) []*common.Location {
	pathTree, ok := m.pathCache.Load(from)
	if ok {
		pathTree := pathTree.(map[common.Point]common.Point)
		pt := to
		path := []*common.Location{}
		for {
			path = append(path, m.Map[pt])
			pt, ok = pathTree[pt]
			if !ok {
				return nil
			}
			if pt == from {
				for i := 0; i < len(path)/2; i++ {
					path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
				}
				return path
			}
		}
	}
	visited := map[common.Point]bool{
		from: true,
	}
	distance := map[common.Point]int{
		from: 0,
	}
	newPathTree := map[common.Point]common.Point{}
	for {
		updated := false
		for pt, loc := range m.Map {
			if _, ok := distance[pt]; !ok {
				distance[pt] = math.MaxInt32
			}
			for _, c := range loc.Connections {
				if !visited[c.Point] {
					continue
				}
				if distance[c.Point]+1 < distance[pt] {
					distance[pt] = distance[c.Point] + 1
					newPathTree[pt] = c.Point
					visited[pt] = true
					updated = true
				}
			}
		}
		if !updated {
			break
		}
	}
	m.pathCache.Store(from, newPathTree)
	return m.Path(from, to)
}

func (m *Map) String(position common.Point) string {
	ret := ""
	for y := 0; y < m.Bounds.Y; y++ {
		for x := 0; x < m.Bounds.X; x++ {
			pt := common.Point{x, y}
			if m.Wall[pt] {
				ret += term.Scolor(85, 85, 85) + "#"
				continue
			}
			if pt == position {
				ret += term.Scolor(85, 255, 255) + "@"
				continue
			}
			if l, ok := m.Map[pt]; ok {
				if l.Trait != "" {
					ret += term.Scolor(255, 85, 255)
				} else {
					ret += term.Scolor(255, 255, 255)
				}
				ret += "."
				continue
			}
			ret += " "
		}
		ret += "\n"
	}
	ret += term.ScolorReset()
	return ret
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <file>", os.Args[0])
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m, err := LoadMap(f)
	if err != nil {
		log.Fatal(err)
	}

	term.Clear()
	term.MoveCursor(1, 1)
	fmt.Print(m.String(m.Start.Point))
	path := m.Path(m.Start.Point, m.Goal.Point)
	term.Color(85, 255, 255)
	for _, l := range append([]*common.Location{m.Start}, path...) {
		term.MoveCursor(l.Point.X+1, l.Point.Y+1)
		print("@")
	}
	term.MoveCursor(1, m.Bounds.Y+2)
	term.ColorReset()
	for name, locations := range m.Portals {
		log.Printf("%s: %d connections", name, len(locations))
	}
	log.Printf("%d steps", len(path))
}
