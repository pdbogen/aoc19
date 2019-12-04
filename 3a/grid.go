package main

import (
	"fmt"
	image "image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type Grid map[Point]map[int]bool

func (g Grid) Bounds() (min, max Point) {
	for pt := range g {
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
	return min, max
}

func (g Grid) String() string {
	min, max := g.Bounds()
	s := ""
	for y := max.Y; y >= min.Y; y-- {
		for x := min.X; x <= max.X; x++ {
			if x == 0 && y == 0 {
				s += "\x1b[1;32m#\x1b[0m"
				continue
			}
			switch len(g[Point{x, y}]) {
			case 0:
				s += "."
			case 1:
				s += "o"
			case 2:
				s += "\x1b[1;31mX\x1b[0m"
			}
		}
		s += "\n"
	}
	return s
}

func (g Grid) Save(file string) error {
	min, max := g.Bounds()
	xShift := -1 * min.X
	yShift := -1 * min.Y
	img := image.NewRGBA(image.Rect(0, 0, max.X+xShift, max.Y+yShift))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)

	for pt, _ := range g {
		img.Set(pt.X+xShift, pt.Y+yShift, color.Black)
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("could not open %q for writing: %v", file, err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("could not encode to PNG: %v", err)
	}
	return nil
}
