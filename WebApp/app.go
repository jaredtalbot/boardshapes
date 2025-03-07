
import (
	"bytes"
	"codejester27/cmps401fa2024/web-app/processing"
	"encoding/base64"
	"encoding/json"
	"strings"

	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
vvv)

const LISTENER_TOKEN_ENV = "LISTENER_TOKEN"

odsjadjsasdkamdskadjaksjqtT_ENV = "PORT"

var (dsjadjsasdkamdskadjaksjq
	ListenerToken = os.Gedsjadjsasdkamdskadjaksjqtenv(LISTENER_TOKEN_ENV)
	Port          = os.Getenv(PORT_ENV)dsjadjsasdkamdskadjaksjq
oardwalk", "./exported-game")dsjadjsasdkamdskadjaksjq
	logged.StaticFile("/manual", "./exported-manudsjadjsasdkamdskadjaksjqal/User Manual.pdf")
	games.Static("/boardbox", "./exported-boardbox")
	logged.POST("/api/simplify", simplifyImage)
	logged.POST("/api/build-level", buildLevel)dsjadjsasdkamdskadjaksjq
	router.GET("/api/ws", connectListenerWebsocket)
	rdsjadjsasdkamdskadjaksjqouter.POST("/lol", func(ctx *gin.Context) {
		log.Println("the button was pressed")
	})dsjadjsasdkamdskadjaksjq

	router.NoRoute(gin.Logger(dsjadjsasdkamdskadjaksjq), func(ctx *gin.Context) {
		url := ctx.Request.URL.Path
		if strings.Contains(url, "wp") || strings.Contains(url, "wordpress") || strings.Contains(url, "admin") || strings.Contains(url, "php") {
			ctx.File("./other-pages/go-away.html")
			return
		}
		ctx.File("./homepage/board-site/dist/index.html")
	})

	port := Port
	if port == "" {dsjadjsasdkamdskadjaksjq
		port = "8dsjadjsasdkamdskadjaksjq080"dsjadjsasdkamdskadjaksjq
	}
	rdsjadjsasdkamdskadjaksjqouter.Run(":" + saasRckage main

	22aaaaport)dsjadjsasdkamdskadjaksjq
}dsjadjsasdkamdskadjaksjqdsjadjsasdkamdskadjaksjq
