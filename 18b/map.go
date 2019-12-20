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
	Starts      []*Location
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
	pathCache = &sync.Map{}
	lengthCache = &sync.Map{}

	ret := &Map{
		Map:   map[Point]*Location{},
		Locks: map[KeySet]*Location{},
		Keys:  map[KeySet]*Location{},
	}

	grid, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Locate starting positions, so that we can also identify sections.
	rows := bytes.Split(grid, []byte{'\n'})
	for y, row := range rows {
		for x, c := range row {
			if c == '@' {
				loc := &Location{
					Point:       Point{x, y},
					Start:       true,
					Section:     len(ret.Starts),
					Connections: map[Direction]*Location{},
				}
				ret.Starts = append(ret.Starts, loc)
				ret.Map[Point{x, y}] = loc
			}
		}
	}

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
			for dir, offset := range map[Direction]Point{
				North: {0, -1},
				South: {0, 1},
				East:  {1, 0},
				West:  {-1, 0},
			} {
				neighbor, ok := ret.Map[Point{x + offset.X, y + offset.Y}]
				if ok {
					neighbor.Connections[dir.Opposite()] = loc
					loc.Connections[dir] = neighbor
				}
			}
			if c == '.' || c == '@' {
				continue
			}
			if c >= 'A' && c <= 'Z' {
				loc.Lock = Keys[c-'A'+'a']
				ret.Locks[Keys[c-'A'+'a']] = loc
			} else if c >= 'a' && c <= 'z' {
				loc.Key = Keys[c]
				ret.Keys[Keys[c]] = loc
			} else {
				return nil, fmt.Errorf("unexpected map character %c", c)
			}
		}
	}

	log.Print("computing sections...")
locations:
	for _, loc := range ret.Map {
		for sectN, start := range ret.Starts {
			if ret.Path(start.Point, loc.Point) != nil {
				loc.Section = sectN
				continue locations
			}
		}
	}

	log.Print("computing constraints...")
	constraints := map[KeySet]KeySet{}
	// for each key, find other keys in the same section that are needed
	for key, kloc := range ret.Keys {
		path := ret.Path(ret.Starts[kloc.Section].Point, kloc.Point)
		for _, pt := range path {
			lock, ok := ret.Map[pt]
			if !ok {
				panic(fmt.Sprintf("no location at %v", pt))
			}
			if lock.Lock == 0 {
				continue
			}
			if ret.Keys[lock.Lock].Section == kloc.Section {
				constraints[key] |= lock.Lock
			}
		}
	}
	ret.Constraints = constraints

	log.Print("computing djikstra trees...")
	wg := &sync.WaitGroup{}
	for _, kloc := range ret.Keys {
		wg.Add(1)
		go func(k, s Point) {
			ret.Path(k, s)
			wg.Done()
		}(kloc.Point, ret.Starts[0].Point)
	}
	wg.Wait()
	log.Print("done.")

	return ret, nil
}

var pathCache = &sync.Map{}

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
	for {
		updated := false
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
					updated = true
				}
			}
		}
		if !updated {
			break
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
// collect the keys given in keys. all keys are assumed to be in the same
// section.
func (m Map) Length(path []byte, keys KeySet) int32 {
	if keys == 0 {
		return 0
	}

	var start Point
	for _, kb := range KeyBytes {
		k := Keys[kb]
		if keys&k > 0 {
			start = m.Starts[m.Keys[k].Section].Point
			break
		}
	}

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

func (m *Map) KeysInSection(section int) KeySet {
	var keys KeySet
	for k, kloc := range m.Keys {
		if kloc.Section == section {
			keys |= k
		}
	}
	return keys
}
