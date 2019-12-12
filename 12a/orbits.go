package main

import (
	"bufio"
	"fmt"
	draw2 "github.com/pdbogen/mapbot/common/draw"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
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

	record := [][]Body{}

	for i := 1; i <= 1000; i++ {
		UpdatePositions(bodies)
		record = append(record, []Body{*bodies[0], *bodies[1], *bodies[2], *bodies[3],})
		log.Printf("%d: %d", i, Energy(bodies))
	}

	minx, maxx, miny, maxy := 0, 0, 0, 0
	for _, entry := range record {
		for _, body := range entry {
			if body.Position.X < minx {
				minx = body.Position.X
			}
			if body.Position.X > maxx {
				maxx = body.Position.X
			}
			if body.Position.Y < miny {
				miny = body.Position.Y
			}
			if body.Position.Y > maxy {
				maxy = body.Position.Y
			}
		}
	}

	img := &gif.GIF{}
	var last []Body
	for _, entry := range record {
		img.Image = append(img.Image, Draw(entry, last, minx, maxx, miny, maxy, 10))
		img.Delay = append(img.Delay, 10)
		img.Disposal = append(img.Disposal, gif.DisposalNone)
		last = entry
	}
	file, _ := os.OpenFile("img.gif", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err := gif.EncodeAll(file, img); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func Draw(system []Body, last []Body, minx, maxx, miny, maxy int, factor int) *image.Paletted {
	img := image.NewPaletted(image.Rect(minx-minx, miny-miny, (maxx-minx)*factor, (maxy-miny)*factor), []color.Color{
		color.Transparent,
		color.Black,
		color.White,
		color.Gray{Y: 64},
	})
	if last == nil {
		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
	} else {
		for _, body := range last {
			x := body.Position.X - minx
			y := body.Position.Y - miny
			draw.Draw(
				img,
				image.Rect(x*factor, y*factor, (x+1)*factor, (y+1)*factor),
				image.Black,
				image.ZP,
				draw.Src)
			draw.Draw(
				img,
				image.Rect(x*factor+2, y*factor+2, (x+1)*factor-2, (y+1)*factor-2),
				image.NewUniform(color.Gray{Y: 64}),
				image.ZP,
				draw.Src)
		}
	}
	for i, body := range system {
		x := body.Position.X - minx
		y := body.Position.Y - miny
		if last != nil {
			lastX := last[i].Position.X - minx
			lastY := last[i].Position.Y - miny
			draw2.Line(img,
				image.Pt(x*factor+factor/2, y*factor+factor/2),
				image.Pt(lastX*factor+factor/2, lastY*factor+factor/2),
				color.Gray{64})
		}
		draw.Draw(img, image.Rect(x*factor, y*factor, (x+1)*factor, (y+1)*factor), image.White, image.ZP, draw.Src)
	}
	return img
}
