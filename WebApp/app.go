package main

import (
	"bytes"
	"codejester27/cmps401fa2024/processing"
	"image"
	"image/png"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func simplifyImage(c *gin.Context) {
	c.SetAccepted("multipart/form-data")

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
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Could not decode your PNG image."})
		return
	}

	img, err = processing.ResizeImage(img)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	imgc := make(chan image.Image)
	go processing.SimplifyImage(img, imgc)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, <-imgc); err != nil {
		panic(err)
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
	c.Header("Content-Disposition", `attachment; filename="simplified-image.png"`)
}

func main() {
	log.Printf("Running with %d CPUs\n", runtime.NumCPU())

	router := gin.Default()

	router.POST("/api/simplify", simplifyImage)

	router.Run("localhost:8080")
}
