package processing

import (
	"image"
	"image/color"
)

const MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION = 150

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

type PixelIndex int
type RegionIndex int

type RegionMap struct {
	regions []*[]PixelIndex
	pixels  map[PixelIndex]RegionIndex
}

func (rm *RegionMap) NewRegion(pixel PixelIndex) (region RegionIndex) {
	region = RegionIndex(len(rm.regions))
	rm.regions = append(rm.regions, &[]PixelIndex{pixel})
	rm.pixels[pixel] = region
	return
}

type PixelImage interface {
	image.Image

	PixOffset(x int, y int) int
}

func (rm *RegionMap) AddPixelToRegion(pixel PixelIndex, region RegionIndex) {
	newPixelArray := append((*rm.regions[region]), pixel)
	rm.regions[region] = &newPixelArray
	rm.pixels[pixel] = region
}

func (rm *RegionMap) GetPixelIndexHasRegion(pixel PixelIndex) (hasRegion bool) {
	_, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionOfPixelIndex(pixel PixelIndex) (regionIndex RegionIndex, hasRegion bool) {
	regionIndex, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionPixels(region RegionIndex) []PixelIndex {
	return *rm.regions[region]
}

func (rm *RegionMap) GetRegions() []*[]PixelIndex {
	return rm.regions
}

func Traverse(img PixelImage, regionMap *RegionMap, px, py int, regionIndex RegionIndex) {

	if pixelIndex := PixelIndex(img.PixOffset(px, py-1)); ColorsBelongInSameRegion(img.At(px, py), img.At(px, py-1)) && !regionMap.GetPixelIndexHasRegion(pixelIndex) {
		regionMap.AddPixelToRegion(pixelIndex, regionIndex)
		Traverse(img, regionMap, px, py-1, regionIndex)
	}
	if pixelIndex := PixelIndex(img.PixOffset(px, py+1)); ColorsBelongInSameRegion(img.At(px, py), img.At(px, py+1)) && !regionMap.GetPixelIndexHasRegion(pixelIndex) {
		regionMap.AddPixelToRegion(pixelIndex, regionIndex)
		Traverse(img, regionMap, px, py+1, regionIndex)
	}
	if pixelIndex := PixelIndex(img.PixOffset(px-1, py)); ColorsBelongInSameRegion(img.At(px, py), img.At(px-1, py)) && !regionMap.GetPixelIndexHasRegion(pixelIndex) {
		regionMap.AddPixelToRegion(pixelIndex, regionIndex)
		Traverse(img, regionMap, px-1, py, regionIndex)
	}
	if pixelIndex := PixelIndex(img.PixOffset(px+1, py)); ColorsBelongInSameRegion(img.At(px, py), img.At(px+1, py)) && !regionMap.GetPixelIndexHasRegion(pixelIndex) {
		regionMap.AddPixelToRegion(pixelIndex, regionIndex)
		Traverse(img, regionMap, px+1, py, regionIndex)
	}
}

func BuildRegionMap(img PixelImage) *RegionMap {
	regionMap := RegionMap{make([]*[]PixelIndex, 0, 20), make(map[PixelIndex]RegionIndex, (img.Bounds().Dx()*img.Bounds().Dy())/4)}

	bd := img.Bounds()

	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			pixelIndex := PixelIndex(img.PixOffset(x, y))

			if img.At(x, y) != White && !regionMap.GetPixelIndexHasRegion(pixelIndex) {
				regionIndex := regionMap.NewRegion(pixelIndex)
				Traverse(img, &regionMap, x, y, regionIndex)
			}
		}
	}

	return &regionMap
}

// hot damn someone rename this function
func ColorsBelongInSameRegion(a color.Color, b color.Color) bool {
	if a == b {
		return true
	} else if (a == Black && b == Red) || (b == Black && a == Red) {
		return true
	} else {
		return false
	}
}

func PixFromOffset(img image.Image, pixIndex PixelIndex) (x int, y int) {
	switch v := img.(type) {
	case *image.Paletted:
		b := v.Bounds()
		x = b.Min.X + (int(pixIndex) % v.Stride)
		y = b.Min.Y + (int(pixIndex) / v.Stride)
	case *image.NRGBA:
		b := v.Bounds()
		x = b.Min.X + (int(pixIndex)%v.Stride)/4
		y = b.Min.Y + (int(pixIndex) / v.Stride)
	}
	return
}

func SimplifyImage(img *image.Image, result chan image.Image) {
	bd := (*img).Bounds()
	newImg := image.NewPaletted(bd, color.Palette{White, Black, Red, Green, Blue})
	// newImg := image.NewNRGBA(bd)

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
		}
	}

	regionMap := BuildRegionMap(newImg)

	// colors := []color.Color{Black, Red, Green, Blue}
	for region := range regionMap.GetRegions() {
		regionPixels := regionMap.GetRegionPixels(RegionIndex(region))
		// randColor := color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
		// randColor := colors[rand.Intn(len(colors))]
		if len(regionPixels) < MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION {
			for _, pixelIndex := range regionPixels {
				x, y := PixFromOffset(newImg, pixelIndex)
				newImg.Set(x, y, White)
			}
		}

		// } else {
		// 	for _, pixelIndex := range regionPixels {
		// 		x, y := PixFromOffset(newImg, pixelIndex)
		// 		newImg.Set(x, y, randColor)
		// 	}
		// }
	}

	result <- newImg
}
