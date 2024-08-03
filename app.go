package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func simplifyImage(img *image.Image, result chan image.Image) {
	bd := (*img).Bounds()
	newImg := image.NewPaletted(bd, color.Palette{White, Black, Red, Green, Blue})
	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			r, g, b, _ := (*img).At(x, y).RGBA()
			r, g, b = r/256, g/256, b/256

			if max(absDiff(r, g), absDiff(g, b), absDiff(r, b)) < 15 {
				// todo: better way to detect black maybe
				if max(r, g, b) > 127 {
					newImg.Set(x, y, White)
				} else {
					newImg.Set(x, y, Black)
				}
			} else if r > g && r > b {
				newImg.Set(x, y, Red)
			} else if g > r && g > b {
				newImg.Set(x, y, Green)
			} else if b > r && b > g {
				newImg.Set(x, y, Blue)
			}
		}
	}

	result <- newImg
}

func simplify(c *gin.Context) {
	c.SetAccepted("image/png")

	fileh, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": `File not found, be sure it is included under the "image" key of your form`})
		return
	}

	if fileh.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Only PNG images are accepted."})
		return
	}

	file, err := fileh.Open()
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	imgc := make(chan image.Image)
	go simplifyImage(&img, imgc)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, <-imgc); err != nil {
		panic(err)
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
	c.Header("Content-Disposition", `attachment; filename="simplified-image.png"`)
}

func main() {
	router := gin.Default()

	router.POST("/api/simplify", simplify)

	router.Run("localhost:8080")
}
