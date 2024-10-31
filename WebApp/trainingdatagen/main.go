package main

import (
	"bytes"
	"codejester27/cmps401fa2024/web-app/processing"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"time"
)

const (
	WHITE_LABEL = 0
	BLACK_LABEL = 1
	RED_LABEL   = 2
	GREEN_LABEL = 3
	BLUE_LABEL  = 4
)

var (
	inputImagePath      = flag.String("i", "input.png", "Path of the non-simplified file used.")
	simplifiedImagePath = flag.String("s", "", "Path of the simplified file used.")

	matrixSize = flag.Int("m", 5, "Set the size of the matrix generated for each pixel.")
)

var Red color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(0), uint8(255)}
var Green color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(0), uint8(255)}
var Blue color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(255), uint8(255)}
var White color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
var Black color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(255)}

var progress, megabytes float64

func oldMatrixDataThing() {
	var buffer bytes.Buffer

	inputFile, err := os.Open(*inputImagePath)
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	bds := img.Bounds()

	go func() {
		fmt.Println()
		for {
			fmt.Printf("\rProgress: %.2f%% | Current Size: %0.3f MB", progress, megabytes)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	buffer.Write(binary.BigEndian.AppendUint16(make([]byte, 0), uint16((*matrixSize)*(*matrixSize))))

	for y := bds.Min.Y; y < bds.Max.Y; y++ {
		for x := bds.Min.X; x < bds.Max.X; x++ {
			px := x - bds.Min.X
			py := y - bds.Min.Y
			progress, megabytes = (float64(px+py*bds.Dx())/float64(bds.Dx()*bds.Dy()))*100.0, float64(buffer.Len())/1_000_000
			for yr := y - *matrixSize/2; yr <= y+*matrixSize/2; yr++ {
				for xr := x - *matrixSize/2; xr <= x+*matrixSize/2; xr++ {
					var r, g, b uint32
					c := img.At(xr, yr)

					if xr < bds.Min.X || yr < bds.Min.Y || xr >= bds.Max.X || yr >= bds.Max.Y {
						r, g, b = uint32(255), uint32(255), uint32(255)
					} else if nrgba, ok := c.(color.NRGBA); ok {
						// use non-alpha-premultiplied colors
						r, g, b = uint32(nrgba.R), uint32(nrgba.G), uint32(nrgba.B)
					} else {
						var a uint32
						// use alpha-premultiplied colors
						r, g, b, a = c.RGBA()
						mult := 65535 / float64(a)
						// undo alpha-premultiplication
						r, g, b = uint32(float64(r)*mult), uint32(float64(g)*mult), uint32(float64(b)*mult)
						// reduce from 0-65535 to 0-255
						r, g, b = r/256, g/256, b/256
					}

					buffer.Write([]byte{byte(r), byte(g), byte(b)})
				}
			}
		}
	}

	output, err := os.Create("output.dat")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	_, err = io.Copy(output, &buffer)
	if err != nil {
		panic(err)
	}
}

func packMatrixData() {
	var pixelDataBuffer bytes.Buffer
	var labelBuffer bytes.Buffer

	generateLabels := *simplifiedImagePath != ""

	inputFile, err := os.Open(*inputImagePath)
	if err != nil {
		panic(err)
	}

	var simplifiedImgFile *os.File
	if generateLabels {
		simplifiedImgFile, err = os.Open(*simplifiedImagePath)
		if err != nil {
			panic(err)
		}
	}

	inputImg, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	var simplifiedImg image.Image
	if generateLabels {
		simplifiedImg, err = png.Decode(simplifiedImgFile)
		if err != nil {
			panic(err)
		}
	}

	bds := inputImg.Bounds()
	if generateLabels && simplifiedImg.Bounds() != bds {
		panic(errors.New("pack-data: non-simplified and simplified images have different bounds"))
	}

	closePrinter := make(chan struct{})
	defer close(closePrinter)

	go func() {
		fmt.Println()
		for {
			select {
			case <-closePrinter:
				return
			default:
				fmt.Printf("\rProgress: %.2f%% | Current Size: %0.3f MB", progress, megabytes)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// stride
	pixelDataBuffer.Write(binary.BigEndian.AppendUint16(make([]byte, 0), uint16(bds.Dx())))

	for y := bds.Min.Y; y < bds.Max.Y; y++ {
		for x := bds.Min.X; x < bds.Max.X; x++ {
			px := x - bds.Min.X
			py := y - bds.Min.Y
			progress, megabytes = (float64(px+py*bds.Dx())/float64(bds.Dx()*bds.Dy()))*100.0, float64(pixelDataBuffer.Len())/1_000_000
			c := processing.GetNRGBA(inputImg.At(x, y))

			pixelDataBuffer.Write([]byte{c.R, c.G, c.B})

			if generateLabels {
				sc := processing.GetNRGBA(simplifiedImg.At(x, y))

				var label uint8
				switch sc {
				case White:
					label = WHITE_LABEL
				case Black:
					label = BLACK_LABEL
				case Red:
					label = RED_LABEL
				case Green:
					label = GREEN_LABEL
				case Blue:
					label = BLUE_LABEL
				default:
					panic(fmt.Errorf("pack-data: simplified image is not simple (pixel with %v)", sc))
				}

				labelBuffer.WriteByte(label)
			}

		}
	}

	pixelDataFile, err := os.Create("pixel_data.dat")
	if err != nil {
		panic(err)
	}
	defer pixelDataFile.Close()

	_, err = io.Copy(pixelDataFile, &pixelDataBuffer)
	if err != nil {
		panic(err)
	}

	if generateLabels {
		labelDataFile, err := os.Create("label_data.dat")
		if err != nil {
			panic(err)
		}
		defer labelDataFile.Close()

		_, err = io.Copy(labelDataFile, &labelBuffer)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	flag.Parse()

	if *matrixSize < 1 {
		panic(errors.New("matrix size must be positive"))
	}

	if *matrixSize%2 != 1 {
		panic(errors.New("matrix size must be odd"))
	}

	packMatrixData()
}
