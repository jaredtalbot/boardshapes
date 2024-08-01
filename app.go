package main

import (
	"bytes"
	"image"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
)

func simplifyImage(img *image.Image, result chan *image.Image) {

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
	}

	file, err := fileh.Open()
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	img, err := png.Decode(file)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
	c.Header("Content-Disposition", `attachment; filename="simplified-image.png"`)
}

func main() {
	router := gin.Default()

	router.POST("/api/simplify", simplify)

	router.Run("localhost:8080")
}
