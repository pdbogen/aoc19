package common

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
	case None:
		return None
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

type Point struct{ X, Y int }

type Location struct {
	Point
	Connections map[Direction]*Location
	Trait       string
}

func (p Point) North() Point {
	return Point{p.X, p.Y - 1}
}

func (p Point) South() Point {
	return Point{p.X, p.Y + 1}
}

func (p Point) East() Point {
	return Point{p.X + 1, p.Y}
}

func (p Point) West() Point {
	return Point{p.X - 1, p.Y}
}
