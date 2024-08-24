package processing

import (
	"image"
	"image/color"
)

func BuildRegionMap(img image.Image) *RegionMap {
	regionMap := RegionMap{make([]*[]Pixel, 0, 20), make(map[Pixel]RegionIndex, (img.Bounds().Dx()*img.Bounds().Dy())/4)}

	bd := img.Bounds()

	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			pixel := Pixel{uint16(x), uint16(y)}

			if img.At(x, y) != White && !regionMap.GetPixelHasRegion(pixel) {
				regionIndex := regionMap.NewRegion(pixel)
				Traverse(img, &regionMap, x, y, regionIndex)
			}
		}
	}

	return &regionMap
}

type Pixel struct {
	X, Y uint16
}
type RegionIndex int

type RegionMap struct {
	regions []*[]Pixel
	pixels  map[Pixel]RegionIndex
}

func (rm *RegionMap) NewRegion(pixel Pixel) (region RegionIndex) {
	region = RegionIndex(len(rm.regions))
	rm.regions = append(rm.regions, &[]Pixel{pixel})
	rm.pixels[pixel] = region
	return
}

func (rm *RegionMap) AddPixelToRegion(pixel Pixel, region RegionIndex) {
	newPixelArray := append((*rm.regions[region]), pixel)
	rm.regions[region] = &newPixelArray
	rm.pixels[pixel] = region
}

func (rm *RegionMap) GetPixelHasRegion(pixel Pixel) (hasRegion bool) {
	_, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionOfPixel(pixel Pixel) (regionIndex RegionIndex, hasRegion bool) {
	regionIndex, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionPixels(region RegionIndex) []Pixel {
	return *rm.regions[region]
}

func (rm *RegionMap) GetRegions() []*[]Pixel {
	return rm.regions
}

func Traverse(img image.Image, regionMap *RegionMap, px, py int, regionIndex RegionIndex) {

	if pixel := (Pixel{uint16(px), uint16(py - 1)}); ColorRegionEquivalence(img.At(px, py), img.At(px, py-1)) && !regionMap.GetPixelHasRegion(pixel) {
		regionMap.AddPixelToRegion(pixel, regionIndex)
		Traverse(img, regionMap, px, py-1, regionIndex)
	}
	if pixel := (Pixel{uint16(px), uint16(py + 1)}); ColorRegionEquivalence(img.At(px, py), img.At(px, py+1)) && !regionMap.GetPixelHasRegion(pixel) {
		regionMap.AddPixelToRegion(pixel, regionIndex)
		Traverse(img, regionMap, px, py+1, regionIndex)
	}
	if pixel := (Pixel{uint16(px - 1), uint16(py)}); ColorRegionEquivalence(img.At(px, py), img.At(px-1, py)) && !regionMap.GetPixelHasRegion(pixel) {
		regionMap.AddPixelToRegion(pixel, regionIndex)
		Traverse(img, regionMap, px-1, py, regionIndex)
	}
	if pixel := (Pixel{uint16(px + 1), uint16(py)}); ColorRegionEquivalence(img.At(px, py), img.At(px+1, py)) && !regionMap.GetPixelHasRegion(pixel) {
		regionMap.AddPixelToRegion(pixel, regionIndex)
		Traverse(img, regionMap, px+1, py, regionIndex)
	}
}

// hot damn someone rename this function
func ColorRegionEquivalence(a color.Color, b color.Color) bool {
	return a == b
}
