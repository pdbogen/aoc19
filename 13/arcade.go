package main

import (
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	"math/rand"
	"time"
)

type Point struct{ X, Y int }
type Pixel struct{ id, lastId int }

func redraw(grid map[Point]*Pixel) {
	fmt.Print("\x1b[")
	//minx, miny, maxx, maxy := 0, 0, 0, 0
	//for pt := range grid {
	//	if pt == (Point{-1, 0}) {
	//		continue
	//	}
	//	if pt.X < minx {
	//		minx = pt.X
	//	}
	//	if pt.X > maxx {
	//		maxx = pt.X
	//	}
	//	if pt.Y < miny {
	//		miny = pt.Y
	//	}
	//	if pt.Y > maxy {
	//		maxy = pt.Y
	//	}
	//}
	blocks := 0
	//for y := miny; y <= maxy; y++ {
	//	for x := minx; x <= maxx; x++ {
	for pt, pixel := range grid {
		if pixel.lastId == pixel.id {
			continue
		}

		if (pt == Point{-1, 0}) {
			fmt.Printf("\x1b[0m\x1B[25;1HScore: % 5d", pixel.id)
			pixel.lastId = pixel.id
			continue
		}
		fmt.Printf("\x1B[%d;%dH", pt.Y+1, pt.X+1)
		lastId := pixel.lastId
		pixel.lastId = pixel.id
		switch pixel.id {
		case 0:
			switch (lastId) {
			case 4:
				fmt.Print("\x1b[1;31m.")
				pixel.lastId = -4
			case -4:
				fmt.Print("\x1b[0;31m.")
				pixel.lastId = -1
			case 2:
				fmt.Print("\x1b[0;32m$")
				pixel.lastId = -2
			case -2:
				fmt.Print("\x1b[0;32ms")
				pixel.lastId = -1
			default:
				fmt.Print(" ")
			case 3:
				fmt.Print("\x1b[0;33m-")
				pixel.lastId = -1
			}
		case 1:
			fmt.Print("\x1b[1;36m#")
		case 2:
			fmt.Printf("\x1b[38;2;%d;%d;%dm$", rand.Intn(128), 127+rand.Intn(128), rand.Intn(128))
			blocks++
		case 3:
			fmt.Print("\x1b[1;33m=")
		case 4:
			fmt.Print("\x1b[1;31m*")
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	print("\x1b[?25l\x1b[2J")
	defer print("\x1b[?25h\x1b[26;1H")
	prog, err := intcode.LoadFile("input")
	if err != nil {
		panic(err)
	}
	prog[0] = 2

	grid := map[Point]*Pixel{}
	input := make(chan int)
	output := make(chan int)
	var ball, paddle Point
	go intcode.Execute(prog, input, output)
	var joystick int
	frames := 0
	for {
		if ball.X < paddle.X {
			joystick = -1
		} else if ball.X > paddle.X {
			joystick = 1
		} else {
			joystick = 0
		}

		var x int
		var ok bool
		select {
		case x, ok = <-output:
		case input <- joystick:
			time.Sleep(2*time.Millisecond)
			redraw(grid)
			frames++
			continue
		}
		if !ok {
			break
		}

		y, ok := <-output
		if !ok {
			break
		}
		id, ok := <-output
		if !ok {
			break
		}
		if id == 4 {
			ball = Point{x, y}
		}
		if id == 3 {
			paddle = Point{x, y}
		}
		if px, ok := grid[Point{x, y}]; ok {
			px.id = id
		} else {
			grid[Point{x, y}] = &Pixel{id, -1}
		}
	}

	for i := 0; i < 3; i++ {
		time.Sleep(5 * time.Millisecond)
		redraw(grid)
	}
}
