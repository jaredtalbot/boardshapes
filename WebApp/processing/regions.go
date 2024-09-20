package processing

import (
	"image"
	"image/color"
	"slices"
)

func BuildRegionMap(img image.Image) *RegionMap {
	regionMap := RegionMap{make([]*Region, 0, 20), make(map[Pixel]RegionIndex, (img.Bounds().Dx()*img.Bounds().Dy())/4)}

	bd := img.Bounds()

	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			pixel := Pixel{uint16(x), uint16(y)}

			if img.At(x, y) != White {
				regionMap.AddPixelToRegionMap(pixel, img)
			}
		}
	}

	regionMap.CleanupRegionMap()
	return &regionMap
}

type Pixel struct {
	X, Y uint16
}

type Vertex struct {
	X, Y uint16
}

type RegionIndex int
type Region []Pixel

type RegionMap struct {
	regions []*Region
	pixels  map[Pixel]RegionIndex
}

func (rm *RegionMap) NewRegion(pixel Pixel) (region RegionIndex) {
	region = RegionIndex(len(rm.regions))
	rm.regions = append(rm.regions, &Region{pixel})
	rm.pixels[pixel] = region
	return
}

func (rm *RegionMap) AddPixelToRegion(pixel Pixel, region RegionIndex) {
	newPixelArray := append((*rm.regions[region]), pixel)
	rm.regions[region] = &newPixelArray
	rm.pixels[pixel] = region
}

func (rm *RegionMap) AddPixelToRegionMap(pixel Pixel, img image.Image) {
	colorP := img.At(int(pixel.X), int(pixel.Y))
	regionLeft, hasRegionLeft := rm.pixels[Pixel{pixel.X - 1, pixel.Y}]
	colorLeft := img.At(int(pixel.X)-1, int(pixel.Y))
	regionAbove, hasRegionAbove := rm.pixels[Pixel{pixel.X, pixel.Y - 1}]
	colorAbove := img.At(int(pixel.X), int(pixel.Y)-1)
	if hasRegionLeft && ColorRegionEquivalence(colorP, colorLeft) {
		if hasRegionAbove && ColorRegionEquivalence(colorP, colorAbove) && regionLeft != regionAbove { // time to merge regions
			pixelsInRegionAbove := rm.regions[regionAbove]
			// grow left region to fit the above region
			mergedRegion := slices.Grow(*rm.regions[regionLeft], len(*rm.regions[regionLeft])+len(*pixelsInRegionAbove)+1)
			// add all pixels in the above region to the left region
			mergedRegion = append(mergedRegion, *pixelsInRegionAbove...)
			rm.regions[regionLeft] = &mergedRegion
			// fix pixel map for all pixels in the above region
			for _, p := range *pixelsInRegionAbove {
				rm.pixels[p] = regionLeft
			}
			rm.regions[regionAbove] = nil
		}
		rm.AddPixelToRegion(pixel, regionLeft)
	} else if hasRegionAbove && ColorRegionEquivalence(colorP, colorAbove) {
		rm.AddPixelToRegion(pixel, regionAbove)
	} else {
		rm.NewRegion(pixel)
	}
}

// cleans up nil regions and rewrites the pixel map
func (rm *RegionMap) CleanupRegionMap() {
	rm.regions = slices.DeleteFunc(rm.regions, func(r *Region) bool { return r == nil })
	pixelsInRegions := 0
	for _, region := range rm.regions {
		pixelsInRegions += len(*region)
	}
	// new fresh map
	rm.pixels = make(map[Pixel]RegionIndex, pixelsInRegions)
	for regionIndex, region := range rm.regions {
		for _, pixel := range *region {
			rm.pixels[pixel] = RegionIndex(regionIndex)
		}
	}
}

func (rm *RegionMap) GetPixelHasRegion(pixel Pixel) (hasRegion bool) {
	_, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegionOfPixel(pixel Pixel) (regionIndex RegionIndex, hasRegion bool) {
	regionIndex, hasRegion = rm.pixels[pixel]
	return
}

func (rm *RegionMap) GetRegion(region RegionIndex) Region {
	if rp := rm.regions[region]; rp != nil {
		return *rp
	}
	return nil
}

func (rm *RegionMap) GetRegions() []*Region {
	return rm.regions
}

func (re *Region) GetBounds() (regionBounds image.Rectangle) {
	regionBounds = image.Rectangle{Min: image.Pt(65535, 65535), Max: image.Pt(0, 0)}
	for _, pixel := range *re {
		if pixel.X < uint16(regionBounds.Min.X) {
			regionBounds.Min.X = int(pixel.X)
		}
		if pixel.Y < uint16(regionBounds.Min.Y) {
			regionBounds.Min.Y = int(pixel.Y)
		}
		if pixel.X+1 > uint16(regionBounds.Max.X) {
			regionBounds.Max.X = int(pixel.X) + 1
		}
		if pixel.Y+1 > uint16(regionBounds.Max.Y) {
			regionBounds.Max.Y = int(pixel.Y) + 1
		}
	}
	return
}

func ColorRegionEquivalence(a color.Color, b color.Color) bool {
	return a == b
}

func FindRegionPosition(region Region) (int, int) {
	corner := region[0]

	for i := 0; i < len(region); i++ {
		if region[i].X < corner.X {
			corner.X = region[i].X
		}
		if region[i].Y < corner.Y {
			corner.Y = region[i].Y
		}
	}

	return int(corner.X), int(corner.Y)
}

func GetColorOfRegion(region Region, img image.Image) color.Color {
	regionColor := img.At(int(region[0].X), int(region[0].Y))
	return regionColor
}
