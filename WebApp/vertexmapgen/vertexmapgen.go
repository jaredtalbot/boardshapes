package main

import (
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
)

var Red color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(0), uint8(255)}
var Green color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(0), uint8(255)}
var Blue color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(255), uint8(255)}
var Black color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(255)}
var White color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
var Cyan color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(255), uint8(255)}
var Magenta color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(255), uint8(255)}
var Yellow color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(0), uint8(255)}
var Blank color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(0)}

type vertex struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	file, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dc := json.NewDecoder(file)
	vs := make([]vertex, 0)

	err = dc.Decode(&vs)
	if err != nil {
		panic(err)
	}

	maxX, maxY := 0, 0
	for _, v := range vs {
		if v.X > maxX-1 {
			maxX = v.X + 1
		}
		if v.Y > maxY-1 {
			maxY = v.Y + 1
		}
	}

	img := image.NewPaletted(image.Rect(0, 0, maxX, maxY), color.Palette{Blank, Red, Green, Blue, Black, White, Cyan, Magenta, Yellow})
	lvs := len(vs)
	for i := range vs {
		a, b := vs[i], vs[(i+1)%lvs]
		switch true {
		case b.X-a.X == -1 && b.Y-a.Y == -1:
			img.SetColorIndex(a.X, a.Y, 1)
		case b.X-a.X == 0 && b.Y-a.Y == -1:
			img.SetColorIndex(a.X, a.Y, 2)
		case b.X-a.X == 1 && b.Y-a.Y == -1:
			img.SetColorIndex(a.X, a.Y, 3)
		case b.X-a.X == -1 && b.Y-a.Y == 0:
			img.SetColorIndex(a.X, a.Y, 4)
		case b.X-a.X == 1 && b.Y-a.Y == 0:
			img.SetColorIndex(a.X, a.Y, 5)
		case b.X-a.X == -1 && b.Y-a.Y == 1:
			img.SetColorIndex(a.X, a.Y, 6)
		case b.X-a.X == 0 && b.Y-a.Y == 1:
			img.SetColorIndex(a.X, a.Y, 7)
		case b.X-a.X == 1 && b.Y-a.Y == 1:
			img.SetColorIndex(a.X, a.Y, 8)
		default:
			panic(errors.New("vertexmapgen: two vertices not adjacent"))
		}
	}

	outf, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	err = png.Encode(outf, img)
	if err != nil {
		panic(err)
	}

}
