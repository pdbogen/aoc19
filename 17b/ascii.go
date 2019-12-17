package main

import (
	"flag"
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"github.com/pdbogen/aoc19/term"
	"log"
)

type Point struct{ X, Y int }

func Print(scaffold map[Point]byte) {
	term.Clear()
	term.MoveCursor(1, 1)
	x, y := 0, 0
	for {
		c, ok := scaffold[Point{x, y}]
		if !ok {
			if x == 0 {
				break
			}
			x = 0
			y++
			fmt.Printf("\n")
			continue
		}
		fmt.Printf("%c", c)
		x++
	}
}

func Intersections(scaffold map[Point]byte) (ret []Point) {
scaffold:
	for p, c := range scaffold {
		if c != 35 {
			continue
		}
		for _, cardinal := range []Point{
			{p.X - 1, p.Y},
			{p.X + 1, p.Y},
			{p.X, p.Y - 1},
			{p.X, p.Y + 1},
		} {
			if scaffold[cardinal] != 35 {
				continue scaffold
			}
		}
		ret = append(ret, p)
	}
	return ret
}

func Bounds(scaffold map[Point]byte) (min Point, max Point) {
	for pt := range scaffold {
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

func Lines(scaffold map[Point]byte) (ret [][2]Point) {
	min, max := Bounds(scaffold)
	for y := min.Y; y <= max.Y; y++ {
		var start *Point = nil
		var end *Point = nil
		for x := min.X; x <= max.X; x++ {
			switch c := scaffold[Point{x, y}]; c {
			case '#':
				fallthrough
			case '^':
				fallthrough
			case 'v':
				fallthrough
			case '<':
				fallthrough
			case '>':
				if start == nil {
					start = &Point{x, y}
				}
				end = &Point{x, y}
			case '.':
				if start != nil && end != nil {
					if *start != *end {
						ret = append(ret, [2]Point{*start, *end})
					}
					start = nil
					end = nil
				}
			default:
				log.Fatalf("unexpected char %c", c)
			}
		}
		if start != nil && end != nil && *start != *end {
			ret = append(ret, [2]Point{*start, *end})
		}
	}
	//return ret
	for x := min.X; x <= max.X; x++ {
		var start *Point = nil
		var end *Point = nil
		for y := min.Y; y <= max.Y; y++ {
			switch c := scaffold[Point{x, y}]; c {
			case '#':
				fallthrough
			case '^':
				fallthrough
			case 'v':
				fallthrough
			case '<':
				fallthrough
			case '>':
				if start == nil {
					start = &Point{x, y}
				}
				end = &Point{x, y}
			case '.':
				if start != nil && end != nil {
					if *start != *end {
						ret = append(ret, [2]Point{*start, *end})
					}
					start = nil
					end = nil
				}
			default:
				log.Fatalf("unexpected char %c", c)
			}
		}
		if start != nil && end != nil && *start != *end {
			ret = append(ret, [2]Point{*start, *end})
		}
	}
	return ret
}

func main() {
	show := flag.Bool("show", false, "if true, show live video")
	flag.Parse()

	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	inCh, outCh := make(chan int), make(chan int)
	go intcode.Execute(prog, inCh, outCh)
	scaffold := map[Point]byte{}
	x := 0
	y := 0
	for c := range outCh {
		if c == 10 {
			x = 0
			y++
			continue
		}
		scaffold[Point{x, y}] = byte(c)
		x++
	}
	lengths := map[int]int{}
	for _, line := range Lines(scaffold) {
		if line[0].X == line[1].X {
			lengths[line[1].Y-line[0].Y+1]++
		} else {
			lengths[line[1].X-line[0].X+1]++
		}
	}
	log.Print(lengths)

	prog[0] = 2
	inCh, outCh = make(chan int), make(chan int)
	go intcode.Execute(prog, inCh, outCh)

	input := [][]byte{
		[]byte("A,B,A,C,B,A,C,B,A,C"),
		[]byte("L,12,L,12,L,6,L,6"),
		[]byte("R,8,R,4,L,12"),
		[]byte("L,12,L,6,R,12,R,8"),
	}

	done := make(chan int)

	if *show {
		input = append(input, []byte("y"))
		go Output(outCh, done)
	} else {
		input = append(input, []byte("n"))
		go func() {
			for c := range outCh {
				if c < 128 {
					fmt.Printf("%c", c)
				} else {
					log.Print(c)
				}
			}
			close(done)
		}()
	}

	for _, line := range input {
		for _, c := range line {
			inCh <- int(c)
		}
		inCh <- 10
	}
	<-done
}

func Output(outCh <-chan int, done chan<- int) {
	term.Clear()
	term.MoveCursor(1, 1)
	term.HideCursor()
	clear := false
	for out := range outCh {
		switch out {
		case 10:
			if clear {
				term.MoveCursor(1, 1)
				clear = false
				continue
			} else {
				clear = true
			}
			fmt.Printf("%c", out)
		case '#':
			term.Color(85, 255, 255)
			print("#")
			clear = false
		case '.':
			term.Color(85, 85, 85)
			print(".")
			clear = false
		case '>':
			fallthrough
		case '<':
			fallthrough
		case '^':
			fallthrough
		case 'v':
			term.Color(255, 85, 85)
			clear = false
			fmt.Printf("%c", out)
		default:
			term.ColorReset()
			term.MoveCursor(45, 1)
			fmt.Printf("%d\n", out)
		}
	}
	close(done)
	term.ShowCursor()
}
