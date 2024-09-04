package processing

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"slices"
	"testing"
)

var Cyan color.NRGBA = color.NRGBA{uint8(0), uint8(255), uint8(255), uint8(255)}
var Magenta color.NRGBA = color.NRGBA{uint8(255), uint8(0), uint8(255), uint8(255)}
var Yellow color.NRGBA = color.NRGBA{uint8(255), uint8(255), uint8(0), uint8(255)}
var Blank color.NRGBA = color.NRGBA{uint8(0), uint8(0), uint8(0), uint8(0)}

func generateRegion(img image.Image) *Region {
	region := make(Region, 0, img.Bounds().Dx()*img.Bounds().Dy())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if c := color.NRGBAModel.Convert(img.At(x, y)); c != White {
				region = append(region, Pixel{uint16(x), uint16(y)})
			}
		}
	}
	return &region
}

func generateMesh(img image.Image) *[]Vertex {
	mesh := make([]Vertex, 0, img.Bounds().Dx()*img.Bounds().Dy())
	var x, y int

Y:
	for y = img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x = img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if img.At(x, y) != Blank {
				break Y
			}
		}
	}

	for {
		mesh = append(mesh, Vertex{uint16(x), uint16(y)})
		c := color.NRGBAModel.Convert(img.At(x, y))
		switch c {
		case Red:
			x--
			y--
		case Green:
			y--
		case Blue:
			x++
			y--
		case Black:
			x--
		case White:
			x++
		case Cyan:
			x--
			y++
		case Magenta:
			y++
		case Yellow:
			x++
			y++
		default:
			panic("your vertex map sucks")
		}

		if mesh[0].X == uint16(x) && mesh[0].Y == uint16(y) {
			break
		}
	}

	return &mesh
}

func generateTestResultImage(name string, mesh *[]Vertex) error {
	file, err := os.Create(fmt.Sprintf("./test_images/test_out_%s.gif", name))
	if err != nil {
		return err
	}

	bounds := image.Rectangle{Min: image.Pt(65535, 65535), Max: image.Pt(0, 0)}
	for _, pixel := range *mesh {
		if pixel.X+1 < uint16(bounds.Min.X) {
			bounds.Min.X = int(pixel.X)
		}
		if pixel.Y < uint16(bounds.Min.Y) {
			bounds.Min.Y = int(pixel.Y)
		}
		if pixel.X+1 > uint16(bounds.Max.X) {
			bounds.Max.X = int(pixel.X) + 1
		}
		if pixel.Y+1 > uint16(bounds.Max.Y) {
			bounds.Max.Y = int(pixel.Y) + 1
		}
	}
	plt := color.Palette{White, Black, Red}
	baseImg := image.NewPaletted(bounds, plt)
	images := make([]*image.Paletted, len(*mesh))
	delays := make([]int, len(*mesh))
	for _, v := range *mesh {
		baseImg.Set(int(v.X), int(v.Y), Black)
	}

	for i, v := range *mesh {
		img := *baseImg
		img.Pix = make([]uint8, len(baseImg.Pix))
		copy(img.Pix, baseImg.Pix)
		img.Set(int(v.X), int(v.Y), Red)
		images[i] = &img
		delays[i] = 10
	}

	g := &gif.GIF{
		Image: images,
		Delay: delays,
		Config: image.Config{
			ColorModel: plt,
			Width:      bounds.Max.X,
			Height:     bounds.Max.Y,
		},
	}

	err = gif.EncodeAll(file, g)
	if err != nil {
		return err
	}

	return file.Close()
}

func TestRegion_CreateMesh(t *testing.T) {
	tests := []struct {
		name     string
		region   *Region
		wantMesh *[]Vertex
	}{
		{
			name: "circle",
		},
		{
			name: "rect",
		},
		{
			name: "head",
		},
	}

	// use supplied test images
	for i, tt := range tests {
		regionf, err := os.Open(fmt.Sprintf("./test_images/test_%s.png", tt.name))
		if err != nil {
			panic(err)
		}

		regionimg, err := png.Decode(regionf)
		if err != nil {
			panic(err)
		}

		tests[i].region = generateRegion(regionimg)

		vertexf, err := os.Open(fmt.Sprintf("./test_images/test_%s_vertexmap.png", tt.name))
		if err != nil {
			panic(err)
		}

		verteximg, err := png.Decode(vertexf)
		if err != nil {
			panic(err)
		}

		tests[i].wantMesh = generateMesh(verteximg)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeshP, err := tt.region.CreateMesh()
			if err != nil {
				t.Fatalf("Error: %s", err)
			}
			gotMesh := *gotMeshP
			wantMesh := *tt.wantMesh

			fail := func() {
				if generateTestResultImage(tt.name, &gotMesh) == nil {
					t.Logf("Generated test result image for %s", tt.name)
				}
				t.Fatalf("Region.CreateMesh() = %v\n want %v", gotMesh, wantMesh)
			}

			firstIndex := slices.Index(wantMesh, gotMesh[0])
			if firstIndex == -1 || len(wantMesh) != len(gotMesh) {
				fail()
			}

			length := len(gotMesh)
			if gotMesh[1] != wantMesh[(firstIndex+1)%length] {
				// second elements don't match, let's try reverse?
				slices.Reverse(wantMesh)
				firstIndex = (length - 1) - firstIndex
				if gotMesh[1] != wantMesh[(firstIndex+1)%length] {
					fail()
				}
			}
			for i := 2; i < length; i++ {
				if gotMesh[i] != wantMesh[(firstIndex+i)%length] {
					fail()
				}
			}
			if generateTestResultImage(tt.name, &gotMesh) == nil {
				t.Logf("Generated test result image for %s", tt.name)
			}
		})
	}
}
