package main

import (
	draw2 "github.com/pdbogen/mapbot/common/draw"
	"image"
	"image/color"
	"image/draw"
)

var red = color.RGBA{R: 255, A: 255}

func Draw(center Point, bounds Point, points map[Point][]Point, target Point) *image.Paletted {
	mul := 20
	ret := image.NewPaletted(image.Rect(0, 0, bounds.X*mul, bounds.Y*mul), []color.Color{
		color.Black,
		color.White,
		red,
	})
	draw.Draw(ret, ret.Bounds(), image.Black, image.ZP, draw.Src)

	for _, pts := range points {
		for _, pt := range pts {
			draw.Draw(
				ret,
				image.Rect(pt.X*mul+2, pt.Y*mul+2, (pt.X+1)*mul-2, (pt.Y+1)*mul-2),
				image.White,
				image.ZP,
				draw.Src,
			)
		}
	}

	draw2.Line(ret, image.Pt(center.X*mul+mul/2, center.Y*mul+mul/2), image.Pt(target.X*mul+mul/2, target.Y*mul+mul/2), red)
	for _, pt := range []Point{center, target} {
		draw.Draw(
			ret,
			image.Rect(pt.X*mul+2, pt.Y*mul+2, (pt.X+1)*mul-2, (pt.Y+1)*mul-2),
			image.NewUniform(red),
			image.ZP,
			draw.Src,
		)
	}
	return ret
}
