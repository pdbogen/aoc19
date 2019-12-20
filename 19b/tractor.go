package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"github.com/pdbogen/aoc19/term"
	"log"
	"math"
	"os"
	"strconv"
)

type Line struct {
	offset int
	width  int
}

func (l Line) String(minOffset int) string {
	ret := term.Scolor(85, 85, 85)
	for i := 0; i < l.offset-minOffset; i++ {
		ret += "."
	}
	ret += term.Scolor(85, 255, 255)
	for i := 0; i < l.width; i++ {
		ret += "#"
	}
	ret += term.ScolorReset()
	return ret
}

const (
	Nothing = 0
	Beam    = 1
)

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	offset := 0
	if len(os.Args) == 2 {
		offset, _ = strconv.Atoi(os.Args[1])
	}

	var lines []Line
	inCh := make(chan int, 2)
	y := offset
	prev := 0
	minOffset := math.MaxInt32
	for {
		l := Line{
			offset: prev,
		}
		x := prev
	scan:
		for {
			if y == 1 || y == 2 {
				break scan
			}
			outCh := make(chan int, 1)
			inCh <- x
			inCh <- y
			intcode.Execute(prog, inCh, outCh)
			switch <-outCh {
			case Nothing:
				if l.width == 0 {
					l.offset++
				} else {
					break scan
				}
			case Beam:
				l.width++
			}
			x++
		}
		if l.offset > 1 {
			prev = l.offset - 1
		} else {
			prev = l.offset
		}
		if l.offset < minOffset {
			minOffset = l.offset
		}
		lines = append(lines, l)
		log.Printf("% 5d: % 5d-% 5d: %s", len(lines)-1+offset, l.offset, l.width, l.String(minOffset))
		y++
		if box(lines, 100) >= 100 {
			last := lines[len(lines)-1]
			log.Printf("top-left is (%d,%d)", last.offset, len(lines)-100+offset)
			break
		}
	}
}

func box(lines []Line, height int) int {
	if len(lines) < height {
		return 0
	}

	maxOffset := 0
	for y := len(lines) - height; y < len(lines); y++ {
		if lines[y].offset > maxOffset {
			maxOffset = lines[y].offset
		}
	}

	minWidth := math.MaxInt32
	for y := len(lines) - height; y < len(lines); y++ {
		width := lines[y].width - maxOffset + lines[y].offset
		if width < minWidth {
			minWidth = width
		}
	}
	return minWidth
}
