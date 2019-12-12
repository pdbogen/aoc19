package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type Point struct {
	X, Y, Z int
}

func (p *Point) Add(addand Point) {
	p.X += addand.X
	p.Y += addand.Y
	p.Z += addand.Z
}

func (p Point) AbsSum() (ret int) {
	if p.X < 0 {
		p.X *= -1
	}
	if p.Y < 0 {
		p.Y *= -1
	}
	if p.Z < 0 {
		p.Z *= -1
	}
	return p.X + p.Y + p.Z
}

type Body struct {
	Position Point
	Velocity Point
}

func (b Body) String() string {
	return fmt.Sprintf("pos=<x=% 3d, y=% 3d, z=% 3d>, vel=<x=% 3d, y=% 3d, z=% 3d>",
		b.Position.X, b.Position.Y, b.Position.Z,
		b.Velocity.X, b.Velocity.Y, b.Velocity.Z)
}

// Adjustment returns the amount by which i changes in order to bring i closer to j
func Adjustment(i, j int) int {
	if i > j {
		return -1
	} else if i < j {
		return 1
	}
	return 0
}

func UpdatePositions(bodies []*Body) {
	for i := 0; i < len(bodies); i++ {
		bodyI := bodies[i]
		for j := i + 1; j < len(bodies); j++ {
			bodyJ := bodies[j]
			bodyI.Velocity.X += Adjustment(bodyI.Position.X, bodyJ.Position.X)
			bodyJ.Velocity.X -= Adjustment(bodyI.Position.X, bodyJ.Position.X)
			bodyI.Velocity.Y += Adjustment(bodyI.Position.Y, bodyJ.Position.Y)
			bodyJ.Velocity.Y -= Adjustment(bodyI.Position.Y, bodyJ.Position.Y)
			bodyI.Velocity.Z += Adjustment(bodyI.Position.Z, bodyJ.Position.Z)
			bodyJ.Velocity.Z -= Adjustment(bodyI.Position.Z, bodyJ.Position.Z)
		}
	}
	for _, body := range bodies {
		body.Position.Add(body.Velocity)
	}
}

func Energy(system []*Body) (ret int) {
	for _, body := range system {
		ret += body.Position.AbsSum() * body.Velocity.AbsSum()
	}
	return ret
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	rdr := bufio.NewReader(f)
	var bodies []*Body
	for {
		b := new(Body)
		_, err := fmt.Fscanf(rdr, "<x=%d, y=%d, z=%d>\n", &b.Position.X, &b.Position.Y, &b.Position.Z)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			panic(err)
		}
		bodies = append(bodies, b)
	}

	record := map[string]map[[8]int]bool{}
	found := map[string]bool{}

	for {
		done := true
		for name, pt := range map[string][8]int{
			"x": {
				bodies[0].Position.X, bodies[1].Position.X, bodies[2].Position.X,bodies[3].Position.X,
				bodies[0].Velocity.X, bodies[1].Velocity.X, bodies[2].Velocity.X,bodies[3].Velocity.X,
			},
			"y": {
				bodies[0].Position.Y, bodies[1].Position.Y, bodies[2].Position.Y,bodies[3].Position.Y,
				bodies[0].Velocity.Y, bodies[1].Velocity.Y, bodies[2].Velocity.Y,bodies[3].Velocity.Y,
			},
			"z": {
				bodies[0].Position.Z, bodies[1].Position.Z, bodies[2].Position.Z,bodies[3].Position.Z,
				bodies[0].Velocity.Z, bodies[1].Velocity.Z, bodies[2].Velocity.Z,bodies[3].Velocity.Z,
			},
		} {
			if found[name] {
				continue
			}
			if record[name] == nil {
				done = false
				record[name] = map[[8]int]bool{
					pt: true,
				}
				continue
			}
			if record[name][pt] {
				found[name] = true
				continue
			}
			record[name][pt] = true
			done = false
		}
		if done {
			break
		}
		UpdatePositions(bodies)
	}

	lcm := 1
	for name, cycle := range record {
		log.Printf("%s LCM(%d,%d)", name, lcm, len(cycle))
		lcm = LCM(lcm, len(cycle))
		log.Println(lcm)
	}

}

func LCM(a, b int) int {
	gcd := a
	k := b
	for k != 0 {
		gcd, k = k, gcd%k
	}
	return a * b / gcd
}
