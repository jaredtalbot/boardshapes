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

	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SimplifyImage(c *gin.Context) {
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

	newImg := processing.SimplifyImage(img, processing.RegionMapOptions{NoColorSeparation: false, AllowWhite: false})

	regionMap := BuildRegionMapForWebAPI(newImg, processing.RegionMapOptions{})

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, newImg); err != nil {
		panic(err)
	}
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	listenerHub.NotifyListeners(ListenerMessage{
		Type: "simplify",
		Attachments: []AttachedFile{
			{
				Name:          fileh.Filename,
				ContentType:   "image/png",
				Base64Content: base64Img,
				Meta: map[string]any{
					"regionCount": len(regionMap.GetRegions()),
				},
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
	c.Data(http.StatusOK, "image/png", buf.Bytes())
	c.Header("Content-Disposition", `attachment; filename="simplified-image.png"`)
}

func CreateShapes(c *gin.Context) {
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
	data := make([]RegionDataV2, 0, numRegions)

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

		points := make([]uint16, len(optimizedShape)*2)
		for i, v := range optimizedShape {
			points[i*2] = v.X
			points[i*2+1] = v.Y
		}

		r := RegionDataV2{i, regionColor, regionColorString, minX, minY, base64Region, points}
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

type RegionDataV2 struct {
	RegionNumber      int         `json:"regionNumber"`
	RegionColor       color.Color `json:"regionColor"`
	RegionColorString string      `json:"regionColorString"`
	CornerX           int         `json:"cornerX"`
	CornerY           int         `json:"cornerY"`
	RegionImage       string      `json:"regionImage"`
	Shape             []uint16    `json:"shape"`
}

type AttachedFile struct {
	Name          string         `json:"name"`
	ContentType   string         `json:"contentType"`
	Base64Content string         `json:"base64Content"`
	Meta          map[string]any `json:"meta,omitempty"`
}

type ListenerMessage struct {
	Type        string         `json:"type"`
	Content     string         `json:"content,omitempty"`
	Attachments []AttachedFile `json:"attachments,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
}

type ListenerHub struct {
	listeners []chan ListenerMessage
	sync.Mutex
}

var listenerHub = ListenerHub{
	listeners: make([]chan ListenerMessage, 0, 5),
}

func (lh *ListenerHub) NotifyListeners(msg ListenerMessage) {
	lh.Lock()
	defer lh.Unlock()
	for _, l := range lh.listeners {
		l <- msg
	}
}

func (lh *ListenerHub) AddListener(listener chan ListenerMessage) {
	lh.Lock()
	defer lh.Unlock()
	lh.listeners = append(lh.listeners, listener)
}

func (lh *ListenerHub) RemoveListener(listener chan ListenerMessage) {
	lh.Lock()
	defer lh.Unlock()
	close(listener)
	lh.listeners = slices.DeleteFunc(lh.listeners, func(c chan ListenerMessage) bool {
		return c == listener
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ConnectListenerWebsocket(c *gin.Context) {
	if c.Query("token") != ListenerToken {
		c.Status(401)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	ch := make(chan ListenerMessage, 1)
	listenerHub.AddListener(ch)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				listenerHub.RemoveListener(ch)
				conn.Close()
				break
			}
		}
	}()

	log.Printf("Listener connected %s", c.Request.RemoteAddr)
	defer log.Printf("Listener disconnected %s", c.Request.RemoteAddr)
	for m := range ch {
		err := conn.WriteJSON(m)
		if err != nil {
			log.Printf("Could not send %T: %s", m, err)
			listenerHub.RemoveListener(ch)
			conn.Close()
			break
		}
	}
}
