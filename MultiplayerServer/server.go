package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const PORT_ENV = "PORT"

var (
	Port = os.Getenv(PORT_ENV)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func giveReadableMessage(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, "<p>This is the multiplayer server for Boardwalk.</p>"+
		"<p>In order to enjoy Boardwalk's multiplayer features, please go to "+
		"<a href=\"https://www.boardmesh.app/boardwalk\">https://www.boardmesh.app/boardwalk</a>.</p>")
}

func checkAvailable(c *gin.Context) {
	c.Status(204)
}

var lobbies = make(map[string]*MultiplayerLobby, 0)

type MultiplayerLobby struct {
	players []MultiplayerPlayer
	sync.RWMutex
}

type MultiplayerPlayer struct {
	Id             uuid.UUID
	MessageChannel chan PlayerStatusMessage
}

type PlayerStatusMessage struct {
	Animation   string   `json:"animation"`
	Frame       int      `json:"frame"`
	Position    Position `json:"position"`
	Name        string   `json:"name"`
	HatId       string   `json:"hatId"`
	HatPosition Position `json:"hatPosition"`
	HatRotation float32  `json:"hatRotation"`
	FacingLeft  bool     `json:"facingLeft"`
	Id          string   `json:"id,omitempty"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func connectMultiplayer(c *gin.Context) {
	lobbyId := c.Query("lobby")
	if lobbyId == "" {
		c.Status(400)
		return
	}

	log.Printf("Player at %s joined lobby %s", c.ClientIP(), lobbyId)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	player := MultiplayerPlayer{id, make(chan PlayerStatusMessage, 5)}
	lobby, ok := lobbies[lobbyId]
	if !ok {
		lobbies[lobbyId] = &MultiplayerLobby{players: make([]MultiplayerPlayer, 0, 1)}
		lobby = lobbies[lobbyId]
	}
	lobby.Lock()
	lobby.players = append(lobby.players, player)
	lobby.Unlock()

	removePlayer := func() {
		lobby.Lock()
		lobby.players = slices.DeleteFunc(lobby.players, func(p MultiplayerPlayer) bool { return player.Id == p.Id })
		lobby.Unlock()
		close(player.MessageChannel)
		conn.Close()
		log.Printf("Player at %s disconnected from lobby %s", c.ClientIP(), lobbyId)
	}

	go func() {
		for {
			_, data, err := conn.ReadMessage()

			if err != nil {
				log.Println(err)
				removePlayer()
				return
			}

			var msg PlayerStatusMessage
			err = json.Unmarshal(data, &msg)

			if err != nil {
				log.Println(err)
				removePlayer()
				return
			}

			msg.Id = id.String()

			lobby.RLock()
			for _, ply := range lobby.players {
				if player.Id == ply.Id {
					continue
				}
				select {
				case ply.MessageChannel <- msg:
				default:
				}
			}
			lobby.RUnlock()
		}
	}()

	for m := range player.MessageChannel {
		err := conn.WriteJSON(m)
		if err != nil {
			log.Println(err)
			removePlayer()
			break
		}
	}
}

func cleanupEmptyLobbies() {
	for {
		lobbiesCleanedUp := 0
		for k, l := range lobbies {
			l.RLock()
			if len(l.players) == 0 {
				delete(lobbies, k)
				lobbiesCleanedUp++
			}
			l.RUnlock()
		}

		if lobbiesCleanedUp > 0 {
			log.Printf("%d multiplayer lobbies cleaned up!\n", lobbiesCleanedUp)
		}

		time.Sleep(60 * time.Second)
	}
}

func main() {
	log.Printf("Running with %d CPUs\n", runtime.NumCPU())

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())

	logged := router.Group("")
	logged.Use(gin.Logger())

	logged.GET("/", giveReadableMessage)
	logged.GET("/check", checkAvailable)
	router.GET("/join", connectMultiplayer)

	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

	go cleanupEmptyLobbies()

	port := Port
	if port == "" {
		port = "443"
	}
	router.RunTLS(":"+port, "multiplayer.boardmesh.app-crt.pem", "multiplayer.boardmesh.app-key.pem")

}
