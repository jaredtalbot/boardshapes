package api

import (
	"codejester27/cmps401fa2024/web-app/processing"
	"errors"
	"image"
	"image/color"
	"os"
)

const PORT_ENV = "PORT"
const LISTENER_TOKEN_ENV = "LISTENER_TOKEN"

var (
	Port          = os.Getenv(PORT_ENV)
	ListenerToken = os.Getenv(LISTENER_TOKEN_ENV)
)

var ErrImageNoSet = errors.New("cannot use an image that has no Set method in buildRegionMapForWebAPI")

type SettableImage = interface {
	image.Image
	Set(x, y int, color color.Color)
}

func BuildRegionMapForWebAPI(img image.Image, options processing.RegionMapOptions) (regionMap *processing.RegionMap) {
	var removedColor color.Color
	if options.AllowWhite {
		removedColor = processing.Blank
	} else {
		removedColor = processing.White
	}

	regionMap = processing.BuildRegionMap(img, options, func(r *processing.Region) bool {
		if len(*r) >= processing.MINIMUM_NUMBER_OF_PIXELS_FOR_VALID_REGION {
			return true
		}
		if i, ok := img.(SettableImage); ok {
			for _, pixel := range *r {
				i.Set(int(pixel.X), int(pixel.Y), removedColor)
			}
		} else {
			panic(ErrImageNoSet)
		}
		return false
	})

	return regionMap
}
