package main

import (
	"bytes"
	"codejester27/cmps401fa2024/processing"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const LISTENER_TOKEN_ENV = "LISTENER_TOKEN"
const PORT_ENV = "PORT"

var (
	ListenerToken = os.Getenv(LISTENER_TOKEN_ENV)
	Port          = os.Getenv(PORT_ENV)
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
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	notifyListeners(ListenerMessage{
		Content: fmt.Sprintf("Simplified image %s successfully.", fileh.Filename),
		Attachments: []AttachedFile{
			{
				Name:          fileh.Filename,
				ContentType:   "image/png",
				Base64Content: base64Img,
			},
		}})
	c.Data(http.StatusOK, "image/png", buf.Bytes())
	c.Header("Content-Disposition", `attachment; filename="simplified-image.png"`)
}

type AttachedFile struct {
	Name          string
	ContentType   string
	Base64Content string
}

type ListenerMessage struct {
	Content     string         `json:"content"`
	Attachments []AttachedFile `json:"attachments"`
}

var listeners = make([]chan ListenerMessage, 0, 5)
var listenersMutex sync.Mutex

func notifyListeners(msg ListenerMessage) {
	listenersMutex.Lock()
	for _, l := range listeners {
		l <- msg
	}
	listenersMutex.Unlock()
}

func addListener(listener chan ListenerMessage) {
	listenersMutex.Lock()
	listeners = append(listeners, listener)
	listenersMutex.Unlock()
}

func removeListener(listener chan ListenerMessage) {
	listenersMutex.Lock()
	close(listener)
	listeners = slices.DeleteFunc(listeners, func(c chan ListenerMessage) bool {
		return c == listener
	})
	listenersMutex.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectWebsocket(c *gin.Context) {
	if c.Query("token") != ListenerToken {
		c.Status(401)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	ch := make(chan ListenerMessage, 1)
	addListener(ch)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				removeListener(ch)
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
			removeListener(ch)
			conn.Close()
			break
		}
	}
}

func main() {
	log.Printf("Running with %d CPUs\n", runtime.NumCPU())

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())

	logged := router.Group("") // I don't like seeing auth tokens in my terminal so we're not logging the websocket requests
	logged.Use(gin.Logger())

	logged.POST("/api/simplify", simplifyImage)
	router.GET("/api/ws", connectWebsocket)

	port := Port
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
