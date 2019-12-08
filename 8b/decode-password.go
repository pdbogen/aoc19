package main

import (
	"github.com/nfnt/resize"
	"github.com/pdbogen/aoc19/sif"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatalf("could not parse input: %v", err)
	}

	inputImage := sif.DecodeImage(string(data), 25, 6)
	inputImage = resize.Resize(640, 640/25*6, inputImage, resize.NearestNeighbor)

	out, err := os.OpenFile("input.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		log.Fatalf("could not open input.png for writing: %v", err)
	}

	if err := png.Encode(out, inputImage); err != nil {
		log.Fatalf("could not encode image to PNG: %v", err)
	}

	if err := out.Close(); err != nil {
		log.Fatalf("error closing image: %v", err)
	}

	out, err = os.OpenFile("input.gif", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		log.Fatalf("could not open input.gif for writing: %v", err)
	}
	inputGif := sif.ToGif(string(data), 25, 6)
	if err := gif.EncodeAll(out, inputGif); err != nil {
		log.Fatalf("could not encode gif: %v", err)
	}
	if err := out.Close(); err != nil {
		log.Fatalf("error closing image.gif: %v", err)
	}
}
