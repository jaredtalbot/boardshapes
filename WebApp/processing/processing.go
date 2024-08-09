package processing

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
)

const MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION = 100

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

type Pixel struct {
	X, Y int
}

type RegionMap struct {
	regions []*[]Pixel
	pixels  map[Pixel]int
}

func NewRegionMap() *RegionMap {
	return &RegionMap{make([]*[]Pixel, 0), make(map[Pixel]int)}
}

func (rm *RegionMap) NewRegion(pixel Pixel) {
	region := len(rm.regions)
	rm.regions = append(rm.regions, &[]Pixel{pixel})
	rm.pixels[pixel] = region
}

func (rm *RegionMap) AddPixelToRegion(pixel Pixel, region int) {
	newPixelArray := append((*rm.regions[region]), pixel)
	rm.regions[region] = &newPixelArray
	rm.pixels[pixel] = region
}

func (rm *RegionMap) GetPixelHasRegion(pixel Pixel) (hasRegion bool) {
	_, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionOfPixel(pixel Pixel) int {
	return rm.pixels[pixel]
}

func (rm *RegionMap) GetRegionPixels(region int) []Pixel {
	return *rm.regions[region]
}

func (rm *RegionMap) GetRegions() []*[]Pixel {
	return rm.regions
}

// hot damn someone rename this function
func ColorsBelongInSameRegion(a color.Color, b color.Color) bool {
	if a == b {
		return true
	} else if (a == Black && b == Red) || (b == Black && a == Red) {
		return true
	} else {
		fmt.Printf("%v != %v\n", a, b)
		return false
	}
}

func SimplifyImage(img *image.Image, result chan image.Image) {
	bd := (*img).Bounds()
	// newImg := image.NewPaletted(bd, color.Palette{White, Black, Red, Green, Blue})
	newImg := image.NewNRGBA(bd)
	regionMap := *NewRegionMap()

	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			r, g, b, _ := (*img).At(x, y).RGBA()
			r, g, b = r/256, g/256, b/256

			var newPixelColor color.NRGBA
			if max(absDiff(r, g), absDiff(g, b), absDiff(r, b)) < 15 {
				// todo: better way to detect black maybe
				if max(r, g, b) > 127 {
					newPixelColor = White
				} else {
					newPixelColor = Black
				}
			} else if r > g && r > b {
				newPixelColor = Red
			} else if g > r && g > b {
				newPixelColor = Green
			} else if b > r && b > g {
				newPixelColor = Blue
			} else {
				newPixelColor = White
			}
			newImg.Set(x, y, newPixelColor)

			if newPixelColor != White {
				// check neighbor pixels
				if x > bd.Min.X && regionMap.GetPixelHasRegion(Pixel{x - 1, y}) && ColorsBelongInSameRegion(newPixelColor, newImg.At(x-1, y)) {
					regionMap.AddPixelToRegion(Pixel{x, y}, regionMap.GetRegionOfPixel(Pixel{x - 1, y}))
				} else if y > bd.Min.Y && regionMap.GetPixelHasRegion(Pixel{x, y - 1}) && ColorsBelongInSameRegion(newPixelColor, newImg.At(x, y-1)) {
					regionMap.AddPixelToRegion(Pixel{x, y}, regionMap.GetRegionOfPixel(Pixel{x, y - 1}))
				} else {
					regionMap.NewRegion(Pixel{x, y})
				}
			}
		}
	}

	for region := range regionMap.GetRegions() {
		regionPixels := regionMap.GetRegionPixels(region)
		randColor := color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
		for _, pix := range regionPixels {
			newImg.Set(pix.X, pix.Y, randColor)
		}
	}

	result <- newImg
}
