package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"sync"
	"sync/atomic"
)

type Map struct {
	Map         map[Point]*Location
	Bounds      Point
	Keys        map[KeySet]*Location
	Locks       map[KeySet]*Location
	Start       *Location
	Constraints map[KeySet]KeySet
}

func (m Map) String() string {
	var ret string
	for y := 0; y <= m.Bounds.Y; y++ {
		for x := 0; x <= m.Bounds.X; x++ {
			if l, ok := m.Map[Point{x, y}]; ok {
				ret += string(l.Symbol())
			} else {
				ret += "#"
			}
		}
		ret += "\n"
	}
	return ret
}

func LoadMap(reader io.Reader) (*Map, error) {
	ret := &Map{
		Map:   map[Point]*Location{},
		Locks: map[KeySet]*Location{},
		Keys:  map[KeySet]*Location{},
	}

	grid, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rows := bytes.Split(grid, []byte{'\n'})
	for y, row := range rows {
		for x, c := range row {
			ret.Bounds = Point{x, y}
			if c == '#' {
				continue
			}
			loc := &Location{
				Point:       Point{x, y},
				Connections: map[Direction]*Location{},
			}
			ret.Map[Point{x, y}] = loc
			west, westOk := ret.Map[Point{x - 1, y}]
			north, northOk := ret.Map[Point{x, y - 1}]
			if westOk {
				west.Connections[East] = loc
				loc.Connections[West] = west
			}
			if northOk {
				north.Connections[South] = loc
				loc.Connections[North] = north
			}
			if c == '.' {
				continue
			}
			if c >= 'A' && c <= 'Z' {
				loc.Lock = Keys[c-'A'+'a']
				ret.Locks[Keys[c-'A'+'a']] = loc
			} else if c >= 'a' && c <= 'z' {
				loc.Key = Keys[c]
				ret.Keys[Keys[c]] = loc
			} else if c == '@' {
				loc.Start = true
				ret.Start = loc
			} else {
				return nil, fmt.Errorf("unexpected map character %c", c)
			}
		}
	}
	constraints := map[KeySet]KeySet{}
	for k, kloc := range ret.Keys {
		pathToKey := ret.Path(ret.Start.Point, kloc.Point)
		for _, pathPt := range pathToKey {
			if ret.Map[pathPt].Lock > 0 {
				constraints[k] |= ret.Map[pathPt].Lock
			}
		}
	}
	ret.Constraints = constraints

	log.Print("pre-computing dijkstra trees...")
	wg := &sync.WaitGroup{}
	for _, k := range ret.Keys {
		wg.Add(1)
		go func(k Point) {
			ret.Path(k, ret.Start.Point)
			wg.Done()
		}(k.Point)
	}
	wg.Wait()
	log.Print("done.")

	return ret, nil
}

var pathCache = sync.Map{}

func (m *Map) Path(from Point, to Point) []Point {
	pathTree, ok := pathCache.Load(from)
	if ok {
		pathTree := pathTree.(map[Point]Point)
		var ret []Point
		pt := to
		for {
			if pt == from {
				break
			}
			var ok bool
			pt, ok = pathTree[pt]
			if !ok {
				return nil
			}
			ret = append([]Point{pt}, ret...)
		}
		return ret
	}
	log.Printf("Computing path tree for %v", from)
	visited := map[Point]bool{
		from: true,
	}
	distances := map[Point]int32{
		from: 0,
	}
	pathTree = map[Point]Point{}
	for len(visited) < len(m.Map) {
		for _, loc := range m.Map {
			if _, ok := distances[loc.Point]; !ok {
				distances[loc.Point] = math.MaxInt32
			}
			for _, conn := range loc.Connections {
				if !visited[conn.Point] {
					continue
				}
				if distances[conn.Point]+1 < distances[loc.Point] {
					distances[loc.Point] = distances[conn.Point] + 1
					pathTree.(map[Point]Point)[loc.Point] = conn.Point
					visited[loc.Point] = true
				}
			}
		}
	}
	pathCache.Store(from, pathTree)
	return m.Path(from, to)
}

func CanTraverse(path []byte, constraints map[KeySet]KeySet) bool {
	var held KeySet
	for _, pathKey := range path {
		req, ok := constraints[Keys[pathKey]]
		if ok && held&req != req {
			return false
		}
		held |= Keys[pathKey]
	}
	return true
}

type lengthCacheKey struct {
	start Point
	keys  KeySet
}

var lengthCache = &sync.Map{}
var lengthHit, lengthMiss int32

// Length calculates the number of steps in the shortest path from start to
// collect the keys given in keys.
func (m Map) Length(path []byte, keys KeySet) int32 {
	if keys == 0 {
		return 0
	}

	start := m.Start.Point
	if len(path) > 0 {
		start = m.Keys[Keys[path[len(path)-1]]].Point
	}

	lck := lengthCacheKey{start, keys}
	v, ok := lengthCache.Load(lck)
	if ok {
		atomic.AddInt32(&lengthHit, 1)
		return v.(int32)
	}
	atomic.AddInt32(&lengthMiss, 1)

	if atomic.LoadInt32(&lengthMiss)%100 == 0 {
		log.Printf("computing length from %v to collect %s (hit rate: %d%%)", string(path), keys,
			(100*atomic.LoadInt32(&lengthHit))/(atomic.LoadInt32(&lengthHit)+atomic.LoadInt32(&lengthMiss)))
	}
	var best int32 = math.MaxInt32
	for _, kb := range KeyBytes {
		k := Keys[kb]
		if keys&k == 0 {
			continue
		}
		keyPt := m.Keys[k].Point

		candidatePath := make([]byte, len(path)+1)
		copy(candidatePath, path)
		candidatePath[len(candidatePath)-1] = KeySymbols[k]

		if !CanTraverse(candidatePath, m.Constraints) {
			continue
		}

		// locations is the marginal path to the proposed next key
		points := m.Path(start, keyPt)
		// length is the best length if we were to choose this key
		length := m.Length(candidatePath, keys^k)

		length += int32(len(points))

		if length < best {
			best = length
		}
	}
	lengthCache.Store(lck, best)
	return best
}

func (m *Map) AllKeys() KeySet {
	var keys KeySet
	for k := range m.Keys {
		keys |= k
	}
	return keys
}
