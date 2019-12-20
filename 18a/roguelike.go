package main

import (
	"fmt"
	"log"
	"os"
)

type Direction int

const (
	None Direction = iota
	North
	South
	East
	West
)

func (d Direction) Opposite() Direction {
	switch (d) {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	}
	return None
}

type Location struct {
	Point       Point
	Key         KeySet
	Lock        KeySet
	Start       bool
	Connections map[Direction]*Location
}

func (l Location) String() (ret string) {
	return fmt.Sprintf("(%d,%d)", l.Point.X, l.Point.Y)
}

func (l Location) Symbol() byte {
	if l.Start {
		return '@'
	}
	if l.Key > 0 {
		return KeySymbols[l.Key]
	}
	if l.Lock > 0 {
		return KeySymbols[l.Lock] - 'a' + 'A'
	}
	return '.'
}

type Point struct {
	X, Y int
}

const givenA = `#########
#b.A.@.a#
#########`

const givenB = `########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`

const givenC = `########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`

const givenD = `#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`

const givenE = `########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	m, err := LoadMap(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(m.Length(nil, m.AllKeys()))
}
