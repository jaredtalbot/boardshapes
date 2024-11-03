package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

var Red color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(0), uint8(255)}
var Green color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(0), uint8(255)}
var Blue color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(255), uint8(255)}
var White color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
var Black color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(255)}

func main() {
	var buff bytes.Buffer
	file, err := os.Open("predicted_labels.dat")
	if err != nil {
		panic(err)
	}
	io.Copy(&buff, file)
	img := image.NewPaletted(image.Rect(0, 0, 1920, 1080), color.Palette{White, Black, Red, Green, Blue})

	for y := 0; y < 1080; y++ {
		for x := 0; x < 1920; x++ {
			label, err := buff.ReadByte()
			if err != nil {
				panic(err)
			}
			img.SetColorIndex(x, y, label)
		}
	}

	outf, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	err = png.Encode(outf, img)
	if err != nil {
		panic(err)
	}
}
