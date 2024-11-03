package main

import (
	"bytes"
	"codejester27/cmps401fa2024/web-app/processing"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"
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
	pixelImagePath      = flag.String("ip", "input.png", "Path of the non-simplified file used for pixel data.")
	labelImagePath      = flag.String("il", "", "Path of the simplified file used for labels.")
	pixelDataOutputPath = flag.String("op", "pixel_data.tpixeldata", "Path to output pixel data.")
	labelDataOutputPath = flag.String("ol", "label_data.tlabeldata", "Path to output label data.")

	matrixSize = flag.Int("m", 5, "Set the size of the matrix generated for each pixel.")
)

var Red color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(0), uint8(255)}
var Green color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(0), uint8(255)}
var Blue color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(255), uint8(255)}
var White color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
var Black color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(255)}

var progress, megabytes float64

func packData(pixelFiles, labelFiles []*os.File) {
	var pixelDataBuffer bytes.Buffer
	var labelBuffer bytes.Buffer

	// matrix
	pixelDataBuffer.Write(binary.BigEndian.AppendUint16(make([]byte, 0), uint16(*matrixSize*(*matrixSize))))

	for i := range pixelFiles {
		var pixelFile, labelFile *os.File
		pixelFile = pixelFiles[i]
		if labelFiles != nil {
			labelFile = labelFiles[i]
		}

		var pixelImg, labelImg image.Image

		pixelImg, err := png.Decode(pixelFile)
		if err != nil {
			panic(err)
		}
		if labelFile != nil {
			labelImg, err = png.Decode(labelFile)
			if err != nil {
				panic(err)
			}
		}

		bds := pixelImg.Bounds()
		if labelImg != nil && labelImg.Bounds() != bds {
			panic(errors.New("pack-data: non-simplified and simplified images have different bounds"))
		}

		closePrinter := make(chan struct{})

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

		for y := bds.Min.Y; y < bds.Max.Y; y++ {
			for x := bds.Min.X; x < bds.Max.X; x++ {
				px := x - bds.Min.X
				py := y - bds.Min.Y
				progress, megabytes = (float64(px+py*bds.Dx())/float64(bds.Dx()*bds.Dy()))*100.0, float64(pixelDataBuffer.Len())/1_000_000
				for yr := y - *matrixSize/2; yr <= y+*matrixSize/2; yr++ {
					for xr := x - *matrixSize/2; xr <= x+*matrixSize/2; xr++ {
						c := processing.GetNRGBA(pixelImg.At(xr, yr))

						pixelDataBuffer.Write([]byte{c.R, c.G, c.B})
					}
				}
				if labelImg != nil {
					sc := processing.GetNRGBA(labelImg.At(x, y))

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
						panic(fmt.Errorf("pack-data: simplified image is not simple (pixel with %v at %d %d)", sc, x, y))
					}

					labelBuffer.WriteByte(label)
				}
			}
		}
		close(closePrinter)
	}

	var compressedDataBuffer bytes.Buffer
	compressor := gzip.NewWriter(&compressedDataBuffer)
	_, err := io.Copy(compressor, &pixelDataBuffer)
	if err != nil {
		panic(err)
	}

	err = compressor.Close()
	if err != nil {
		panic(err)
	}

	pixelDataFile, err := os.Create(*pixelDataOutputPath)
	if err != nil {
		panic(err)
	}
	defer pixelDataFile.Close()

	_, err = io.Copy(pixelDataFile, &compressedDataBuffer)
	if err != nil {
		panic(err)
	}

	if labelFiles != nil {
		labelDataFile, err := os.Create(*labelDataOutputPath)
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

func getFiles() (pixelFiles, labelFiles []*os.File) {

	if pixelDir, err := os.ReadDir(*pixelImagePath); err == nil {
		pixelDir = slices.DeleteFunc(pixelDir, func(e fs.DirEntry) bool { return e.IsDir() || e.Name() == strings.TrimSuffix(e.Name(), ".png") })
		if len(pixelDir) == 0 {
			log.Fatalf("%s has no files", *pixelImagePath)
		}

		pixelFileNames := make([]string, len(pixelDir))
		for i, v := range pixelDir {
			pixelFileNames[i] = v.Name()
		}

		if *labelImagePath != "" {
			labelDir, err := os.ReadDir(*pixelImagePath)
			if err != nil {
				panic(err)
			}
			labelDir = slices.DeleteFunc(labelDir, func(e fs.DirEntry) bool { return e.IsDir() || e.Name() == strings.TrimSuffix(e.Name(), ".png") })
			if len(pixelDir) != len(labelDir) {
				log.Fatalln("pixel and label dirs have inequal length")
			}

			labelFileNames := make([]string, len(labelDir))
			for i, v := range labelDir {
				labelFileNames[i] = v.Name()
			}

			slices.Sort(labelFileNames)

			labelFiles = make([]*os.File, len(labelFileNames))
			for i, v := range labelFileNames {
				labelFiles[i], err = os.Open(*labelImagePath + "/" + v)
				if err != nil {
					panic(err)
				}
			}
		}

		slices.Sort(pixelFileNames)
		pixelFiles = make([]*os.File, len(pixelFileNames))
		for i, v := range pixelFileNames {
			pixelFiles[i], err = os.Open(*pixelImagePath + "/" + v)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	pixelFile, err := os.Open(*pixelImagePath)
	if err != nil {
		panic(err)
	}

	pixelFiles = []*os.File{pixelFile}

	if *labelImagePath != "" {
		labelFile, err := os.Open(*labelImagePath)
		if err != nil {
			panic(err)
		}
		labelFiles = []*os.File{labelFile}
	}
	return
}

func main() {
	flag.Parse()

	if *matrixSize < 1 {
		panic(errors.New("matrix size must be positive"))
	}

	if *matrixSize%2 != 1 {
		panic(errors.New("matrix size must be odd"))
	}

	pixelFiles, labelFiles := getFiles()

	packData(pixelFiles, labelFiles)
}
