package api

import (
	"bytes"
	"codejester27/cmps401fa2024/web-app/processing"
	"encoding/base64"
	"encoding/json"

	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegionDataV1 struct {
	RegionNumber      int                 `json:"regionNumber"`
	RegionColor       color.Color         `json:"regionColor"`
	RegionColorString string              `json:"regionColorString"`
	CornerX           int                 `json:"cornerX"`
	CornerY           int                 `json:"cornerY"`
	RegionImage       string              `json:"regionImage"`
	Mesh              []processing.Vertex `json:"mesh"`
}

func BuildLevel(c *gin.Context) {
	c.SetAccepted("multipart/form-data")

	fileh, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": `File not found, be sure it is included under the "image" key of your form`})
		return
	}

	file, err := fileh.Open()
	if err != nil {
		panic(err)
	}

	preserveColor := c.Request.FormValue("preserveColor") == "true"
	noColorSeparation := c.Request.FormValue("noColorSeparation") == "true"
	allowWhite := c.Request.FormValue("allowWhite") == "true"

	contentType := fileh.Header.Get("Content-Type")
	var img image.Image
	switch contentType {
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Could not decode your PNG image."})
			return
		}
	case "image/jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Could not decode your JPEG image."})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Only JPEG and PNG images are accepted."})
		return
	}

	img, err = processing.ResizeImage(img)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	opts := processing.RegionMapOptions{
		NoColorSeparation: noColorSeparation,
		AllowWhite:        allowWhite,
	}

	newImg := processing.SimplifyImage(img, opts)

	regionMap := BuildRegionMapForWebAPI(newImg, opts)

	regions := regionMap.GetRegions()
	numRegions := len(regions)
	data := make([]RegionDataV1, 0, numRegions)

	for i := range numRegions {
		region := regionMap.GetRegionByIndex(i)

		minX, minY := processing.FindRegionPosition(region)
		regionColor := processing.GetColorOfRegion(region, newImg, noColorSeparation)
		var regionColorString string

		switch regionColor {
		case processing.Red:
			regionColorString = "Red"
		case processing.Green:
			regionColorString = "Green"
		case processing.Blue:
			regionColorString = "Blue"
		case processing.Black:
			regionColorString = "Black"
		case processing.White:
			regionColorString = "White"
		}

		regionImage := image.NewNRGBA(region.GetBounds())

		if preserveColor {
			for j := 0; j < len(*region); j++ {
				regionImage.Set(int((*region)[j].X), int((*region)[j].Y), img.At(int((*region)[j].X), int((*region)[j].Y)))
			}
		} else {
			for j := 0; j < len(*region); j++ {
				regionImage.Set(int((*region)[j].X), int((*region)[j].Y), regionColor)
			}
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, regionImage); err != nil {
			panic(err)
		}
		base64Region := base64.StdEncoding.EncodeToString(buf.Bytes())

		shape, err := region.CreateShape()
		if err != nil {
			continue
		}
		optimizedShape := processing.StraightOpt(shape)
		r := RegionDataV1{i, regionColor, regionColorString, minX, minY, base64Region, optimizedShape}
		data = append(data, r)
	}

	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	listenerHub.NotifyListeners(ListenerMessage{
		Type: "build-level",
		Attachments: []AttachedFile{
			{
				Name:          "level.json",
				ContentType:   "application/json",
				Base64Content: string(d),
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
	c.Data(http.StatusOK, "application/json", d)
}
