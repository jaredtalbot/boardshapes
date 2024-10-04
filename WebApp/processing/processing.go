package processing

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
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

func manhattanDistance(a Vertex, b Vertex) int {
	return absDiff(int(a.X), int(b.X)) + absDiff(int(a.Y), int(b.Y))
}

var ErrImageTooWide = errors.New("image is too wide")

type RegionPixelStatus uint8

type RegionPixel struct {
	InRegion, Visited bool
}

func (region *Region) CreateMesh() (mesh *[]Vertex, err error) {
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
		regionPixels[int(v.X+1)-regionBounds.Min.X][int(v.Y+1)-regionBounds.Min.Y].InRegion = true
	}

	OuterVertexMesh := make([]Vertex, 0, regionBounds.Dx()*2+regionBounds.Dy()*2-2)

	vertexesToVisit := []Vertex{{0, 0}}
	// visit outer pixel
	for len(vertexesToVisit) > 0 {
		v := vertexesToVisit[len(vertexesToVisit)-1]
		vertexesToVisit = vertexesToVisit[:len(vertexesToVisit)-1]
		if !regionPixels[v.X][v.Y].Visited {
			regionPixels[v.X][v.Y].Visited = true
			if v.X > 0 && !regionPixels[v.X-1][v.Y].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X - 1, v.Y}) {
				if regionPixels[v.X-1][v.Y].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X - 1, v.Y})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X - 1, v.Y})
				}
			}
			if v.X < uint16(len(regionPixels))-1 && !regionPixels[v.X+1][v.Y].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X + 1, v.Y}) {
				if regionPixels[v.X+1][v.Y].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X + 1, v.Y})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X + 1, v.Y})
				}
			}
			if v.Y > 0 && !regionPixels[v.X][v.Y-1].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X, v.Y - 1}) {
				if regionPixels[v.X][v.Y-1].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X, v.Y - 1})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X, v.Y - 1})
				}
			}
			if v.Y < uint16(len(regionPixels[0]))-1 && !regionPixels[v.X][v.Y+1].Visited && !slices.Contains(OuterVertexMesh, Vertex{v.X, v.Y + 1}) {
				if regionPixels[v.X][v.Y+1].InRegion {
					OuterVertexMesh = append(OuterVertexMesh, Vertex{v.X, v.Y + 1})
				} else {
					vertexesToVisit = append(vertexesToVisit, Vertex{v.X, v.Y + 1})
				}
			}
		}
	}

	vertexMatrix := make([][]bool, regionBounds.Dx())
	for i := range vertexMatrix {
		vertexMatrix[i] = make([]bool, regionBounds.Dy())
	}

	// translate all vertices by (-1, -1)
	// necessary because we added extra space for the region up above
	for i, v := range OuterVertexMesh {
		OuterVertexMesh[i].X--
		OuterVertexMesh[i].Y--
		vertexMatrix[v.X-1][v.Y-1] = true
	}

	wither(vertexMatrix)

	OuterVertexMesh = slices.DeleteFunc(OuterVertexMesh, func(v Vertex) bool {
		return !vertexMatrix[v.X][v.Y]
	})

	var previousVertex Vertex
	var isPreviousVertexSet = false
	var currentVertex Vertex = OuterVertexMesh[0]
	SortedOuterVertexMesh := make([]Vertex, 0, len(OuterVertexMesh))

	for {
		adjacentVertices := make([]Vertex, 0, 8)

		forAdjacents(currentVertex.X, currentVertex.Y, len(vertexMatrix), len(vertexMatrix[0]), func(x, y uint16) {
			if vertexMatrix[x][y] {
				adjacentVertices = append(adjacentVertices, Vertex{uint16(x), uint16(y)})
			}
		})

		// sort by manhattan distance to put diagonal vertices last
		slices.SortFunc(adjacentVertices, func(a Vertex, b Vertex) int {
			return manhattanDistance(a, currentVertex) - manhattanDistance(b, currentVertex)
		})

		if !isPreviousVertexSet {
			isPreviousVertexSet = true
			previousVertex = adjacentVertices[0]
			SortedOuterVertexMesh = append(SortedOuterVertexMesh, previousVertex)
		}

		SortedOuterVertexMesh = append(SortedOuterVertexMesh, currentVertex)

		// scary!!!
		if len(adjacentVertices) == 3 {
			var a, b *Vertex = nil, nil
			for _, v := range adjacentVertices {
				if v != previousVertex {
					var adjs uint8 = 0
					forAdjacents(v.X, v.Y, len(vertexMatrix), len(vertexMatrix[0]), func(x, y uint16) {
						if vertexMatrix[x][y] {
							adjs++
						}
					})
					if adjs == 2 {
						if a == nil {
							a = &v
						} else {
							return nil, errors.New("region-to-mesh: mesh generation failed")
						}
					} else if adjs == 3 {
						if b == nil {
							b = &v
						} else {
							return nil, errors.New("region-to-mesh: mesh generation failed")
						}
					}
				}
			}

			if a == nil || b == nil {
				return nil, errors.New("region-to-mesh: mesh generation failed")
			}

			vertexMatrix[a.X][a.Y] = false
			adjacentVertices = slices.DeleteFunc(adjacentVertices, func(v Vertex) bool {
				return v == *a
			})

			wither(vertexMatrix)
		}

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

		if currentVertex == SortedOuterVertexMesh[0] {
			return &SortedOuterVertexMesh, nil
		}

		if len(SortedOuterVertexMesh) >= len(OuterVertexMesh) {
			return nil, errors.New("region-to-mesh: could not close mesh")
		}
	}
}

func wither(matrix [][]bool) {
	maxX, maxY := len(matrix), len(matrix[0])
	frames := make([]*image.Paletted, 0)
	for {
		frames = append(frames, MatrixToImage(matrix))
		verticesToRemove := make([]Vertex, 0)
		for x := uint16(0); x < uint16(maxX); x++ {
		Y:
			for y := uint16(0); y < uint16(maxY); y++ {
				if matrix[x][y] {
					// add adjacents to slice
					adjacentVertices := make([]Vertex, 0, 8)
					forAdjacents(x, y, maxX, maxY, func(x, y uint16) {
						if matrix[x][y] {
							adjacentVertices = append(adjacentVertices, Vertex{x, y})
						}
					})

					if len(adjacentVertices) < 2 {
						verticesToRemove = append(verticesToRemove, Vertex{x, y})
						continue Y
					}

					for _, v := range adjacentVertices {
						noAdjacents := true
						for _, ov := range adjacentVertices {
							if v != ov {
								vX, vY, ovX, ovY := int(v.X), int(v.Y), int(ov.X), int(ov.Y)
								noAdjacents = vX < ovX-1 || vX > ovX+1 || vY < ovY-1 || vY > ovY+1
								if !noAdjacents {
									break
								}
							}
						}
						if noAdjacents {
							continue Y
						}
					}

					verticesToRemove = append(verticesToRemove, Vertex{x, y})
				}
			}
		}

		if len(verticesToRemove) < 1 {
			break
		}

		for _, v := range verticesToRemove {
			matrix[v.X][v.Y] = false
		}
	}
	frames = append(frames, MatrixToImage(matrix))
	delays := make([]int, len(frames))
	for i := range frames {
		delays[i] = 15
	}

	g := &gif.GIF{
		Image: frames,
		Delay: delays,
		Config: image.Config{
			ColorModel: color.Palette{White, Black},
			Width:      maxX,
			Height:     maxY,
		},
	}

	gifFile, err := os.Create("withering.gif")

	if err != nil {
		panic(err)
	}

	err = gif.EncodeAll(gifFile, g)

	if err != nil {
		panic(err)
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
