package processing

import (
	"errors"
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
)

const MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION = 50

var Red color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(0), uint8(255)}
var Green color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(0), uint8(255)}
var Blue color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(255), uint8(255)}
var White color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
var Black color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(255)}

func absDiff[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a T, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}

var ErrImageTooWide = errors.New("image is too wide")

func ResizeImage(img image.Image) (image.Image, error) {
	const MAX_HEIGHT = 1080
	const MAX_WIDTH = 2000

	bd := img.Bounds()
	if bd.Dy() > MAX_HEIGHT {
		scalar := float64(MAX_HEIGHT) / float64(bd.Dy())
		newWidth := math.Round(float64(bd.Dx()) * scalar)
		if newWidth > MAX_WIDTH {
			return nil, ErrImageTooWide
		}
		scaledImg := image.NewNRGBA(image.Rect(0, 0, int(newWidth), MAX_HEIGHT))
		draw.NearestNeighbor.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
		return scaledImg, nil
	} else if bd.Dx() > MAX_WIDTH {
		return nil, ErrImageTooWide
	}
	return img, nil
}

func SimplifyImage(img image.Image) (result image.Image, regionCount int) {
	bd := img.Bounds()
	newImg := image.NewPaletted(bd, color.Palette{White, Black, Red, Green, Blue})
	// newImg := image.NewNRGBA(bd)

	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r/256, g/256, b/256

			var newPixelColor color.NRGBA
			avg := (r + g + b) / 3
			if max(absDiff(avg, r), absDiff(avg, g), absDiff(avg, b)) < 10 {
				// todo: better way to detect black maybe
				if max(r, g, b) > 115 {
					newPixelColor = White
				} else {
					newPixelColor = Black
				}
			} else if r > g && r > b {
				newPixelColor = Red
			} else if g > r && (g > b || b-g < 10) {
				newPixelColor = Green
			} else if b > r && b > g {
				newPixelColor = Blue
			} else {
				newPixelColor = White
			}
			newImg.Set(x, y, newPixelColor)
		}
	}

	regionMap := BuildRegionMap(newImg)

	// colors := []color.Color{Black, Red, Green, Blue}
	for region := range regionMap.GetRegions() {
		region := regionMap.GetRegion(RegionIndex(region))
		// randColor := color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
		// randColor := colors[rand.Intn(len(colors))]
		if len(region) < MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION {
			for _, pixel := range region {
				newImg.Set(int(pixel.X), int(pixel.Y), White)
			}
		} else {
			regionCount++
		}

		// } else {
		// 	for _, pixel := range regionPixels {
		// 		newImg.Set(int(pixel.X), int(pixel.Y), randColor)
		// 	}
		// }
	}

	return newImg, regionCount
}
