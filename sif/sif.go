package sif

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
)

const (
	Black       = 0
	White       = 1
	Transparent = 2
)

var index = map[int]color.Color{
	0: color.Black,
	1: color.White,
	2: color.Transparent,
}

func FlattenLayers(layers [][][]int) (image [][]int) {
	ret := make([][]int, len(layers))
	for x, col := range layers {
		ret[x] = make([]int, len(col))
		for y, row := range col {
			for _, pixel := range row {
				if pixel == Transparent {
					continue
				}
				ret[x][y] = pixel
				break
			}
		}
	}
	return ret
}

func DecodeImage(data string, width, height int) image.Image {
	ret := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(ret, ret.Bounds(), image.NewUniform(color.Transparent), image.ZP, draw.Src)

	img := FlattenLayers(ToLayers(data, width, height))
	for x, col := range img {
		for y, pixel := range col {
			ret.Set(x, y, index[pixel])
		}
	}
	return ret
}

// ToLayers returns pixels[x][y][layer]
func ToLayers(data string, width, height int) [][][]int {
	layers := len(data) / (width * height)
	ret := make([][][]int, width)
	for x := 0; x < width; x++ {
		ret[x] = make([][]int, height)
		for y := 0; y < height; y++ {
			ret[x][y] = make([]int, layers)
			for l := 0; l < layers; l++ {
				dataIdx := l*width*height + y*width + x
				ret[x][y][l] = int(data[dataIdx] - '0')
			}
		}
	}
	return ret
}

func ToGif(data string, width, height int) *gif.GIF {
	mul := 640 / width
	if width*mul < 640 {
		mul++
	}

	layers := ToLayers(data, width, height)
	nLayers := len(data) / width / height
	ret := &gif.GIF{
		Image:     make([]*image.Paletted, nLayers),
		Delay:     make([]int, nLayers),
		Disposal:  make([]byte, nLayers),
	}

	for i := 0; i < nLayers; i++ {
		ret.Image[i] = image.NewPaletted(
			image.Rect(0, 0, width*mul, height*mul),
			[]color.Color{index[0], index[1], index[2]},
		)
		draw.Draw(ret.Image[i], ret.Image[i].Bounds(), image.Transparent, image.ZP, draw.Src)
		ret.Delay[i] = 10
		ret.Disposal[i] = gif.DisposalNone
	}
	ret.Delay[nLayers-1] = 1000

	for x, col := range layers {
		for y, stack := range col {
			for l := range stack {
				//for l := nLayers - 1; l >= 0; l-- {
				//ret.Image[l].Set(x, y, index[stack[l]])
				draw.Draw(
					ret.Image[nLayers-1-l], image.Rect(x*mul, y*mul, (x+1)*mul, (y+1)*mul),
					image.NewUniform(index[stack[l]]),
					image.ZP,
					draw.Src,
				)
				//for ix := x * mul; ix < x*mul+mul; ix++ {
				//	for iy := y * mul; iy < y*mul+mul; iy++ {
				//		ret.Image[l].Set(ix, iy, index[stack[l]])
				//	}
				//}
			}
		}
	}
	return ret
}
