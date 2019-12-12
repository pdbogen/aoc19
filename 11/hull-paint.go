package main

import (
	"github.com/pdbogen/aoc19/intcode"
	image "image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const (
	black = 0
	white = 1

	Up    = 0
	Right = 1
	Down  = 2
	Left  = 3
)

type Point struct{ X, Y int }

var hull = map[Point]int{}

func main() {
	program, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatalf("could not load program: %v", err)
	}

	inCh := make(chan int, 1)
	outCh := make(chan int)
	done := make(chan int)
	go func() {
		if _, err := intcode.Execute(program, inCh, outCh); err != nil {
			log.Fatalf("running program: %v", err)
		}
		close(done)
	}()

	robot := Point{}
	dir := Up
	hull[robot] = 1
loop:
	for {
		select {
		case inCh <- hull[robot]:
		case <-done:
			break loop
		}

		hull[robot] = <-outCh
		log.Printf("painted %v -> %v", robot, hull[robot])

		if <-outCh == 0 {
			// left
			if dir == Up {
				dir = Left
			} else {
				dir -= 1
			}
		} else {
			if dir == Left {
				dir = Up
			} else {
				dir += 1
			}
		}
		switch dir {
		case Up:
			robot.Y--
		case Right:
			robot.X++
		case Down:
			robot.Y++
		case Left:
			robot.X--
		}
	}
	var min, max Point
	for pt := range hull {
		if min.X > pt.X {
			min.X = pt.X
		}
		if min.Y > pt.Y {
			min.Y = pt.Y
		}
		if max.X < pt.X {
			max.X = pt.X
		}
		if max.Y < pt.Y {
			max.Y = pt.Y
		}
	}
	log.Printf("Bounds %v - %v", min, max)
	log.Printf("# points: %d", len(hull))

	img := image.NewRGBA(image.Rect(0, 0, max.X*10, max.Y*10))
	draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
	for pt, color := range hull {
		if color == white {
			draw.Draw(img, image.Rect(pt.X*10, pt.Y*10, pt.X*10+10, pt.Y*10+10), image.White, image.ZP, draw.Src)
		}
	}
	f, _ := os.OpenFile("output.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	png.Encode(f, img)
}
