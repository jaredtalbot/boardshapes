package processing

import (
	"errors"
	"image"
	"image/color"
	"math"
	"slices"

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

type RegionPixelStatus uint8

type RegionPixel struct {
	InRegion, Visited bool
}

func (region *Region) CreateMesh() (mesh *[]Vertex) {
	if len(*region) == 0 {
		// TODO: handle bad case
	}
	regionBounds := image.Rectangle{Min: image.Pt(65535, 65535), Max: image.Pt(0, 0)}
	for _, pixel := range *region {
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

	// will sastisfy my requirements.
	VertexMesh := make([][]RegionPixel, regionBounds.Dx()+2)
	for i := range VertexMesh {
		VertexMesh[i] = make([]RegionPixel, regionBounds.Dy()+2)
	}

	for _, v := range *region {
		VertexMesh[int(v.X+1)-regionBounds.Min.X][int(v.Y+1)-regionBounds.Min.Y].InRegion = true
	}

	OuterVertexMesh := make([]Vertex, 0, 4)

	vertexesToVisit := []Vertex{{0, 0}}
	// visit outer pixel
	for len(vertexesToVisit) > 0 {
		v := vertexesToVisit[len(vertexesToVisit)-1]
		vertexesToVisit = vertexesToVisit[:len(vertexesToVisit)-1]
		if !VertexMesh[v.X][v.Y].Visited {
			VertexMesh[v.X][v.Y].Visited = true
			if v.X > 0 && !VertexMesh[v.X-1][v.Y].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X - 1, v.Y}) {
				if VertexMesh[v.X-1][v.Y].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X - 1, v.Y})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X - 1, v.Y})
				}
			}
			if v.X < uint16(len(VertexMesh))-1 && !VertexMesh[v.X+1][v.Y].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X + 1, v.Y}) {
				if VertexMesh[v.X+1][v.Y].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X + 1, v.Y})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X + 1, v.Y})
				}
			}
			if v.Y > 0 && !VertexMesh[v.X][v.Y-1].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X, v.Y - 1}) {
				if VertexMesh[v.X][v.Y-1].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X, v.Y - 1})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X, v.Y - 1})
				}
			}
			if v.Y < uint16(len(VertexMesh[0]))-1 && !VertexMesh[v.X][v.Y+1].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X, v.Y + 1}) {
				if VertexMesh[v.X][v.Y+1].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X, v.Y + 1})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X, v.Y + 1})
				}
			}
		}
	}

	// Note: It may be worth getting rid of the OuterVertexMesh slice entirely, and instead
	//		 just adding an "OuterVertex" boolean to the RegionPixel struct so that you
	//		 can just reuse the VertexMesh variable.
	//		 This would make checking adjacent pixels much easier when creating a sorted outer mesh.

	// sort outermesh

	SortedOuterVertexMesh := make([]Vertex, 0, len(OuterVertexMesh))
	tempVertexStack := make([]Vertex, 0, len(OuterVertexMesh))
	onStack := false
	for i := 0; i < len(OuterVertexMesh); i++ {
		currentVertex := OuterVertexMesh[i]

		if len(tempVertexStack) > 0 && onStack {

			for j := 0; j < len(tempVertexStack); j++ {
				if j+1 < len(tempVertexStack) && onStack {
					nextVertexTemp := tempVertexStack[j+1]

					if currentVertex.X-nextVertexTemp.X < 1 {
						SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertexTemp)
						currentVertex = nextVertexTemp
						onStack = false

					} else if currentVertex.Y-nextVertexTemp.Y < 1 {
						SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertexTemp)
						currentVertex = nextVertexTemp
						onStack = false

					} else if currentVertex.X < nextVertexTemp.X && currentVertex.Y-nextVertexTemp.Y < 1 {
						SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertexTemp)
						currentVertex = nextVertexTemp
						onStack = false

					} else if currentVertex.X > nextVertexTemp.X && nextVertexTemp.Y-currentVertex.Y < 1 {
						SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertexTemp)
						currentVertex = nextVertexTemp
						onStack = false
					}
				}
			}
		}
		if i+1 < len(OuterVertexMesh) {
			nextVertex := OuterVertexMesh[i+1]

			if currentVertex.X-nextVertex.X < 1 {
				SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertex)
				currentVertex = nextVertex

			} else if currentVertex.Y-nextVertex.Y < 1 {
				SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertex)
				currentVertex = nextVertex

			} else if currentVertex.X < nextVertex.X && currentVertex.Y-nextVertex.Y < 1 {
				SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertex)
				currentVertex = nextVertex

			} else if currentVertex.X > nextVertex.X && nextVertex.Y-currentVertex.Y < 1 {
				SortedOuterVertexMesh = append(SortedOuterVertexMesh, nextVertex)
				currentVertex = nextVertex

			} else {
				tempVertexStack = append(tempVertexStack, nextVertex)
				currentVertex = nextVertex
				onStack = true

			}
		}
	}
	return &SortedOuterVertexMesh
}

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
