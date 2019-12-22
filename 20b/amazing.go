package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/pdbogen/aoc19/common"
	"github.com/pdbogen/aoc19/term"
	draw2 "github.com/pdbogen/mapbot/common/draw"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
)

type Map struct {
	Map       map[common.Point]*common.Location
	Wall      map[common.Point]bool
	Bounds    common.Point
	Portals   map[string][]*common.Location
	Start     *common.Location
	Goal      *common.Location
	pathCache *sync.Map
	Layers    int
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

func (m *Map) Path(from, to common.Point, renderAnim bool) []*common.Location {
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
	anim := &gif.GIF{}
	log.Printf("%d: running djikstra with %d points", m.Layers, len(m.Map))
	distance := make([][]int, m.Bounds.X)
	n := 0
	for x := range distance {
		distance[x] = make([]int, m.Bounds.Y)
		for y := range distance[x] {
			n++
			distance[x][y] = math.MaxInt32
		}
	}
	log.Printf("%d: initialized %d cells", m.Layers, n)
	distance[from.X][from.Y] = 0
	queue := []*common.Location{m.Map[from]}
	newPathTree := map[common.Point]common.Point{}
	for {
		var nextQueue = []*common.Location{}
		for _, loc := range queue {
			pt := loc.Point
			for _, c := range loc.Connections {
				cPt := c.Point
				if distance[cPt.X][cPt.Y] > distance[pt.X][pt.Y]+1 {
					distance[cPt.X][cPt.Y] = distance[pt.X][pt.Y] + 1
					nextQueue = append(nextQueue, m.Map[cPt])
					newPathTree[cPt] = pt
				}
			}
		}
		if distance[to.X][to.Y] != math.MaxInt32 || len(nextQueue) == 0 {
			break
		}
		queue = nextQueue
		if renderAnim {
			anim.Image = append(anim.Image, m.ImageDjikstra(distance))
			anim.Disposal = append(anim.Disposal, gif.DisposalNone)
			anim.Delay = append(anim.Delay, 10)
		}
	}
	m.pathCache.Store(from, newPathTree)
	path := m.Path(from, to, false)
	if renderAnim {
		log.Printf("%d: constructing a GIF...", m.Layers)
		if path != nil {
			anim.Image = append(anim.Image, m.Image(path...))
			anim.Disposal = append(anim.Disposal, gif.DisposalNone)
			anim.Delay = append(anim.Delay, 500)
		}
		var err error
		for i := range anim.Image {
			if i == 0 {
				continue
			}
			anim.Image[i], err = DiffOnly(anim.Image[i-1], anim.Image[i])
			if err != nil {
				panic(err)
			}
		}
		f, err := os.OpenFile("map.gif", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
		if err != nil {
			panic(err)
		}
		if err := gif.EncodeAll(f, anim); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
	return path
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

func (m *Map) Layer(n int) (layered *Map) {
	columns := int(math.Sqrt(float64(n)))
	rows := int(math.Ceil(float64(n) / float64(columns)))

	layered = &Map{
		Map:       map[common.Point]*common.Location{},
		Wall:      map[common.Point]bool{},
		Bounds:    common.Point{m.Bounds.X * columns, m.Bounds.Y * rows},
		Portals:   map[string][]*common.Location{},
		pathCache: &sync.Map{},
		Layers:    n,
	}

	layerPortals := map[int]map[string]*common.Location{}

	for l := 0; l < n; l++ {
		offsetX := (m.Bounds.X - 1) * (l % columns)
		offsetY := (m.Bounds.Y - 1) * (l / columns)
		layerPortals[l] = map[string]*common.Location{}
		for pt, loc := range m.Map {
			portal := loc.Trait != "" && loc.Trait != "AA" && loc.Trait != "ZZ"
			goal := loc.Trait == "AA" || loc.Trait == "ZZ"
			outer := portal && (
				pt.X == 2 || pt.X == m.Bounds.X-2 ||
					pt.Y == 2 || pt.Y == m.Bounds.Y-2)

			newPt := common.Point{pt.X + offsetX, pt.Y + offsetY}

			// No outer portals on layer 0
			if l == 0 && outer {
				layered.Wall[newPt] = true
				continue
			}

			// No goals below layer 0
			if l > 0 && goal {
				layered.Wall[newPt] = true
				continue
			}

			// Ok! Create a new location for this layer, with X appropriate X shift.
			newLoc := &common.Location{
				Point:       newPt,
				Connections: map[common.Direction]*common.Location{},
				Trait:       loc.Trait,
			}
			layered.Map[newPt] = newLoc

			if newLoc.Trait == "AA" {
				layered.Start = newLoc
			}
			if newLoc.Trait == "ZZ" {
				layered.Goal = newLoc
			}

			if portal && !outer {
				layerPortals[l][loc.Trait] = newLoc
			}
			if outer {
				// the inner portal on the higher layer must already exist
				newLoc.Connections[common.None] = layerPortals[l-1][loc.Trait]
				layerPortals[l-1][loc.Trait].Connections[common.None] = newLoc
			}

			// Note that we iterate through the original map in some random order!
			for d, p := range map[common.Direction]common.Point{
				common.North: newPt.North(),
				common.East:  newPt.East(),
				common.South: newPt.South(),
				common.West:  newPt.West(),
			} {
				if connectedLoc, ok := layered.Map[p]; ok {
					newLoc.Connections[d] = connectedLoc
					connectedLoc.Connections[d.Opposite()] = newLoc
				}
			}
		}
		for pt := range m.Wall {
			newPt := common.Point{pt.X + offsetX, pt.Y + offsetY}
			layered.Wall[newPt] = true
		}
	}

	return layered
}

func (m *Map) ImageDjikstra(distances [][]int) *image.Paletted {
	points := []common.Point{}
	for x := range distances {
		for y, d := range distances[x] {
			if d != math.MaxInt32 {
				points = append(points, common.Point{x, y})
			}
		}
	}
	return m.ImagePoints(points)
}

func (m *Map) Image(path ...*common.Location) *image.Paletted {
	points := []common.Point{}
	for _, p := range path {
		points = append(points, p.Point)
	}
	return m.ImagePoints(points)
}

func (m *Map) ImagePoints(points []common.Point) *image.Paletted {
	factor := 5
	img := image.NewPaletted(image.Rect(0, 0, m.Bounds.X*factor+factor, m.Bounds.Y*factor+factor), []color.Color{
		color.Black,
		color.White,
		color.RGBA{255, 85, 85, 255},
		color.RGBA{85, 255, 85, 255},
		color.RGBA{85, 255, 255, 255},
		color.RGBA{85, 85, 85, 255},
		color.RGBA{255, 85, 255, 255},
		color.Transparent,
	})

	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	magenta := map[common.Point]bool{}
	for _, pt := range points {
		magenta[pt] = true
		draw.Draw(
			img,
			image.Rect((pt.X*factor)+1, (pt.Y*factor)+1, (pt.X*factor)+(factor-2), (pt.Y*factor)+(factor-2)),
			image.NewUniform(color.RGBA{255, 85, 255, 255}),
			image.Point{},
			draw.Src)
	}

	for pt, loc := range m.Map {
		var c color.Color = color.White
		if loc.Trait == "AA" || loc.Trait == "ZZ" {
			c = color.RGBA{85, 255, 85, 255}
		} else if magenta[loc.Point] {
			c = color.RGBA{255, 85, 255, 255}
		} else if loc.Trait != "" {
			c = color.RGBA{85, 255, 255, 255}
		}
		draw.Draw(
			img,
			image.Rect((pt.X*factor)+1, (pt.Y*factor)+1, (pt.X*factor)+(factor-2), (pt.Y*factor)+(factor-2)),
			image.NewUniform(c),
			image.Point{},
			draw.Src)
	}

	for pt := range m.Wall {
		draw.Draw(
			img,
			image.Rect((pt.X*factor)+1, (pt.Y*factor)+1, (pt.X*factor)+(factor-2), (pt.Y*factor)+(factor-2)),
			image.NewUniform(color.RGBA{85, 85, 85, 255}),
			image.Point{},
			draw.Src)
	}

	for _, loc := range m.Map {
		for dir, conLoc := range loc.Connections {
			if dir != common.None {
				continue
			}
			c := color.RGBA{255, 85, 85, 255}
			if magenta[loc.Point] && magenta[conLoc.Point] {
				c = color.RGBA{255, 85, 255, 255}
			}

			draw2.Line(img,
				image.Pt(loc.Point.X*factor+factor/2-1, loc.Point.Y*factor+factor/2-1),
				image.Pt(conLoc.Point.X*factor+factor/2-1, conLoc.Point.Y*factor+factor/2-1),
				c)
		}
	}

	return img
}

func Worker(m *Map, in <-chan int, out chan<- struct {
	m *Map
	p []*common.Location
}, running *int32, anim bool) {
	for n := range in {
		log.Printf("%d: constructing map...", n)
		m := m.Layer(n)
		log.Printf("%d: calculating path...", n)
		path := m.Path(m.Start.Point, m.Goal.Point, anim)
		if path != nil {
			out <- struct {
				m *Map
				p []*common.Location
			}{m: m, p: path}
		}
	}
	atomic.AddInt32(running, -1)
	out <- struct {
		m *Map
		p []*common.Location
	}{m: nil, p: nil}
}

func main() {
	anim := flag.Bool("anim", false, "if set, produce an animation of Djikstra at work")
	mapfile := flag.String("file", "input", "file from which to load our map")
	start := flag.Int("start", 20, "minimum layer count to attempt")
	num := flag.Int("num", -1, "how many values to try; -1 means try forever")
	flag.Parse()

	f, err := os.Open(*mapfile)
	if err != nil {
		log.Fatal(err)
	}
	m, err := LoadMap(f)
	if err != nil {
		log.Fatal(err)
	}

	maps := make(chan int)
	results := make(chan struct {
		m *Map
		p []*common.Location
	})

	var running int32 = 4
	for i := 0; i < 4; i++ {
		go Worker(m, maps, results, &running, *anim)
	}

	n := *start - 1
	var bestMap *Map
	var bestPath []*common.Location
	best := math.MaxInt32
layerLoop:
	for {
		n++
		if *num >= 0 && (n-*start) >= *num && maps != nil {
			close(maps)
			maps = nil
		}
		select {
		case maps <- n:
		case result := <-results:
			if result.m == nil {
				if atomic.LoadInt32(&running) == 0 {
					break layerLoop
				}
				continue
			}

			path := result.p
			if path == nil {
				log.Printf("%d: No Path", result.m.Layers)
				continue
			}

			if len(path) < best {
				log.Printf("%d: path length %d is better than %d", n, len(path), best)
				best = len(path)
				bestMap = result.m
				bestPath = path
				continue
			}

			log.Printf("%d: path length did not improve! I guess we are done.", n)
			break layerLoop
		}
	}

	if bestMap == nil {
		log.Fatal("I guess not. Sorry, kids.")
	}

	bestImage := bestMap.Image(bestPath...)
	f, err = os.OpenFile("map.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, bestImage); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func DiffOnly(a, b *image.Paletted) (*image.Paletted, error) {
	if a.Bounds() != b.Bounds() {
		return nil, errors.New("bounds mismatch")
	}
	if !reflect.DeepEqual(a.Palette, b.Palette) {
		return nil, errors.New("palette mismatch")
	}
	ret := image.NewPaletted(a.Bounds(), a.Palette)
	draw.Draw(ret, ret.Bounds(), image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	for x := 0; x < a.Bounds().Max.X; x++ {
		for y := 0; y < a.Bounds().Max.Y; y++ {
			bc := b.At(x, y)
			if a.At(x, y) == bc {
				continue
			}
			ret.Set(x, y, bc)
		}
	}
	return ret, nil
}
