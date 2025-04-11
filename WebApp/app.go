package main

import (
	"codejester27/cmps401fa2024/web-app/api"
	"strings"

	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func containsAny(haystack string, needles ...string) bool {
	for _, s := range needles {
		if strings.Contains(haystack, s) {
			return true
		}
	}
	return false
}

func main() {
	log.Printf("Running with %d CPUs\n", runtime.NumCPU())

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())

	logged := router.Group("") // I don't like seeing auth tokens in my terminal so we're not logging the websocket requests
	logged.Use(gin.Logger())

	// cors
	logged.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

	games := logged.Group("")

	games.Use(func(ctx *gin.Context) {
		ctx.Header("Cross-Origin-Embedder-Policy", "require-corp")
		ctx.Header("Cross-Origin-Opener-Policy", "same-origin")
	})

	logged.StaticFile("/", "./homepage/board-site/dist/index.html")
	logged.StaticFile("/board.png", "./homepage/board-site/dist/board.png")
	logged.Static("/assets", "./homepage/board-site/dist/assets")
	games.Static("/boardwalk", "./exported-game")
	logged.StaticFile("/manual", "./exported-manual/User Manual.pdf")
	games.Static("/boardbox", "./exported-boardbox")
	logged.POST("/api/simplify", api.SimplifyImage)
	logged.POST("/api/build-level", api.BuildLevel)
	logged.POST("/api/create-shapes", api.CreateShapes)
	router.GET("/api/ws", api.ConnectListenerWebsocket)

	router.NoRoute(gin.Logger(), func(ctx *gin.Context) {
		url := ctx.Request.URL.Path
		if containsAny(url, "wp", "wordpress", "admin", "php") {
			ctx.File("./misc-assets/teapot.txt")
			ctx.Status(http.StatusTeapot)
			return
		}
		ctx.File("./homepage/board-site/dist/index.html")
	})

	port := api.Port
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
