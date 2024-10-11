package processing

import (
	"cmp"
	"errors"
	"fmt"
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

// func manhattanDistance(a Vertex, b Vertex) int {
// 	return absDiff(int(a.X), int(b.X)) + absDiff(int(a.Y), int(b.Y))
// }

var ErrImageTooWide = errors.New("image is too wide")

type RegionPixelStatus uint8

type RegionPixel byte

const (
	REGION_PIXEL_IN_REGION = 0b00000001
	REGION_PIXEL_VISITED   = 0b00000010
	REGION_PIXEL_IS_OUTER  = 0b00000100
	REGION_PIXEL_IS_INNER  = 0b00001000
)

func (r *RegionPixel) MarkInRegion() {
	*r = *r | REGION_PIXEL_IN_REGION
}

func (r *RegionPixel) MarkVisited() {
	*r = *r | REGION_PIXEL_VISITED
}

func (r *RegionPixel) MarkIsOuter() {
	*r = *r | REGION_PIXEL_IS_OUTER
}

func (r *RegionPixel) MarkIsInner() {
	*r = *r | REGION_PIXEL_IS_INNER
}

func (r RegionPixel) InRegion() bool {
	return r&REGION_PIXEL_IN_REGION > 0
}

func (r RegionPixel) Visited() bool {
	return r&REGION_PIXEL_VISITED > 0
}

func (r RegionPixel) IsOuter() bool {
	return r&REGION_PIXEL_IS_OUTER > 0
}

func (r RegionPixel) IsInner() bool {
	return r&REGION_PIXEL_IS_INNER > 0
}

func (r RegionPixel) String() string {
	return fmt.Sprintf("in region: %t; visited: %t; in mesh: %t", r.InRegion(), r.Visited(), r.IsOuter())
}

func (region *Region) CreateMesh() (mesh []Vertex, err error) {
	if len(*region) == 0 {
		return nil, errors.New("region-to-mesh: region is empty")
	}
	regionBounds := region.GetBounds()

	// will sastisfy my requirements.
	regionPixels := make([][]RegionPixel, regionBounds.Dx()+2)
	for i := range regionPixels {
		regionPixels[i] = make([]RegionPixel, regionBounds.Dy()+2)
	}

	for _, v := range *region {
		regionPixels[int(v.X+1)-regionBounds.Min.X][int(v.Y+1)-regionBounds.Min.Y].MarkInRegion()
	}

	verticesToVisit := []Vertex{{0, 0}}
	// visit outer pixels
	for len(verticesToVisit) > 0 {
		v := verticesToVisit[len(verticesToVisit)-1]
		verticesToVisit = verticesToVisit[:len(verticesToVisit)-1]
		if !regionPixels[v.X][v.Y].Visited() {
			regionPixels[v.X][v.Y].MarkVisited()
			forNonDiagonalAdjacents(v.X, v.Y, len(regionPixels), len(regionPixels[0]), func(x, y uint16) {
				if !regionPixels[x][y].Visited() && !regionPixels[x][y].IsOuter() {
					if regionPixels[x][y].InRegion() {
						regionPixels[x][y].MarkIsOuter()
					} else {
						verticesToVisit = append(verticesToVisit, Vertex{x, y})
					}
				}
			})

		}
	}

	vertexMeshes := make([][]Vertex, 0, 1)

	// find inner pixels
	for y := uint16(0); y < uint16(len(regionPixels[0])); y++ {
		for x := uint16(0); x < uint16(len(regionPixels)); x++ {
			rp := regionPixels[x][y]
			// check if inner pixel
			if !rp.Visited() && !rp.IsOuter() {
				verticesToVisit := []Vertex{{x, y}}
				newInnerMesh := make([]Vertex, 0, regionBounds.Dx()+regionBounds.Dy())
				// visit inner pixels
				for len(verticesToVisit) > 0 {
					v := verticesToVisit[len(verticesToVisit)-1]
					verticesToVisit = verticesToVisit[:len(verticesToVisit)-1]
					if !regionPixels[v.X][v.Y].Visited() {
						regionPixels[v.X][v.Y].MarkVisited()
						forNonDiagonalAdjacents(v.X, v.Y, len(regionPixels), len(regionPixels[0]), func(x, y uint16) {
							if !regionPixels[x][y].Visited() && !regionPixels[x][y].IsInner() {
								if regionPixels[x][y].IsOuter() {
									regionPixels[x][y].MarkIsInner()
									newInnerMesh = append(newInnerMesh, Vertex{x, y})
								} else {
									verticesToVisit = append(verticesToVisit, Vertex{x, y})
								}
							}
						})

					}
				}
				vertexMeshes = append(vertexMeshes, newInnerMesh)
			}
		}
	}

	vertexMatrix := make([][]bool, regionBounds.Dx())
	for i := range vertexMatrix {
		vertexMatrix[i] = make([]bool, regionBounds.Dy())
	}

	if len(vertexMeshes) == 0 {
		return nil, errors.New("region-to-mesh: region is too thin")
	}

	vertexMesh := slices.MaxFunc(vertexMeshes, func(a, b []Vertex) int {
		return cmp.Compare(len(a), len(b))
	})

	// translate all vertices by (-1, -1)
	// necessary because we added extra space for the region up above
	for i, v := range vertexMesh {
		vertexMesh[i].X--
		vertexMesh[i].Y--
		vertexMatrix[v.X-1][v.Y-1] = true
	}

	var previousVertex Vertex
	var isPreviousVertexSet = false
	var currentVertex Vertex = vertexMesh[0]
	sortedOuterVertexMesh := make([]Vertex, 0, len(vertexMesh))

	for {
		adjacentVertices := make([]Vertex, 0, 8)

		forAdjacents(currentVertex.X, currentVertex.Y, len(vertexMatrix), len(vertexMatrix[0]), func(x, y uint16) {
			if vertexMatrix[x][y] {
				adjacentVertices = append(adjacentVertices, Vertex{uint16(x), uint16(y)})
			}
		})

		if !isPreviousVertexSet {
			isPreviousVertexSet = true
			previousVertex = adjacentVertices[0]
			sortedOuterVertexMesh = append(sortedOuterVertexMesh, previousVertex)
		}

		sortedOuterVertexMesh = append(sortedOuterVertexMesh, currentVertex)

		if len(adjacentVertices) != 2 {
			return nil, errors.New("region-to-mesh: mesh generation failed")
		}

		if adjacentVertices[0] == previousVertex {
			previousVertex = currentVertex
			currentVertex = adjacentVertices[1]
		} else {
			previousVertex = currentVertex
			currentVertex = adjacentVertices[0]
		}

		if currentVertex == sortedOuterVertexMesh[0] {
			return sortedOuterVertexMesh, nil
		}

		if len(sortedOuterVertexMesh) >= len(vertexMesh) {
			return nil, errors.New("region-to-mesh: could not close mesh")
		}
	}
}

func PrintMatrix(matrix [][]bool) {
	for _, s := range matrix {
		for _, v := range s {
			if v {
				fmt.Print("██")
			} else {
				fmt.Print("░░")
			}
		}
		fmt.Println()
	}
}

func MatrixToImage(matrix [][]bool) *image.Paletted {
	maxX, maxY := len(matrix), len(matrix[0])
	result := image.NewPaletted(image.Rect(0, 0, maxX, maxY), color.Palette{White, Black})
	for x := uint16(0); x < uint16(maxX); x++ {
		for y := uint16(0); y < uint16(maxY); y++ {
			if matrix[x][y] {
				result.SetColorIndex(int(x), int(y), 1)
			}
		}
	}
	return result
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
			// for _, pixel := range region {
			// 	newImg.Set(int(pixel.X), int(pixel.Y), Black)
			// }
		}
	}

	return newImg, regionCount
}

func forNonDiagonalAdjacents(x, y uint16, maxX, maxY int, function func(x, y uint16)) {
	if y > 0 {
		function(x, y-1)
	}
	if x > 0 {
		function(x-1, y)
	}
	if x < uint16(maxX)-1 {
		function(x+1, y)
	}
	if y < uint16(maxY)-1 {
		function(x, y+1)
	}
}

func forAdjacents(x, y uint16, maxX, maxY int, function func(x, y uint16)) {
	if y > 0 {
		if x > 0 {
			function(x-1, y-1)
		}
		function(x, y-1)
		if x < uint16(maxX)-1 {
			function(x+1, y-1)
		}
	}
	if x > 0 {
		function(x-1, y)
	}
	if x < uint16(maxX)-1 {
		function(x+1, y)
	}
	if y < uint16(maxY)-1 {
		if x > 0 {
			function(x-1, y+1)
		}
		function(x, y+1)
		if x < uint16(maxX)-1 {
			function(x+1, y+1)
		}
	}
}
