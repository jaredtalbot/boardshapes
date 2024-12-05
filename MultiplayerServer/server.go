package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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

var lobbies = make(map[string]*MultiplayerLobby, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		originUrl, err := url.Parse(origin)
		if err != nil {
			return false
		}
		hostname := originUrl.Hostname()

		switch hostname {
		case "cmps401fa2024.onrender.com", "www.boardmesh.app", "boardmesh.app", "multiplayer.boardmesh.app", "localhost":
			return true
		}
		return false
	},
}

type MultiplayerLobby struct {
	Id      uuid.UUID
	Name    string
	Players []MultiplayerPlayer
	sync.RWMutex
}

func (ml *MultiplayerLobby) ToSummaryDto() MultiplayerLobbySummaryDto {
	return MultiplayerLobbySummaryDto{Id: ml.Id.String(), Name: ml.Name}
}

type MultiplayerLobbySummaryDto struct {
	Id, Name string
}

type MultiplayerPlayer struct {
	Id             uuid.UUID
	Name           string
	IsHost         bool
	MessageChannel chan PlayerMessage
}

func (mpp MultiplayerPlayer) ToDto() MultiplayerPlayerDto {
	return MultiplayerPlayerDto{Id: mpp.Id.String(), Name: mpp.Name}
}

type MultiplayerPlayerDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlayerMessage struct {
	MsgType string `json:"type"`
	Content any    `json:"content,omitempty"`
}

type PlayerPosition struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
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

func getLobbies(c *gin.Context) {
	lobbyDtos := make([]MultiplayerLobbySummaryDto, 0, len(lobbies))
	for _, v := range lobbies {
		lobbyDtos = append(lobbyDtos, v.ToSummaryDto())
	}
	c.JSON(200, lobbyDtos)
}

func connectMultiplayer(c *gin.Context) {
	lobbyId := c.Query("lobby")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	createdLobby := false
	var lobby *MultiplayerLobby
	if lobbyId == "" {
		lobbyId = uuid.New().String()
		lobbies[lobbyId] = &MultiplayerLobby{
			Name:    "New Lobby",
			Players: make([]MultiplayerPlayer, 0, 1),
		}
		lobby = lobbies[lobbyId]
		createdLobby = true
	} else {
		var ok bool
		lobby, ok = lobbies[lobbyId]
		if !ok {
			conn.Close()
			c.Status(400)
			return
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	player := MultiplayerPlayer{
		Id:             id,
		Name:           "Player",
		IsHost:         createdLobby,
		MessageChannel: make(chan PlayerMessage, 5),
	}

	lobby.Lock()
	lobby.Players = append(lobby.Players, player)
	lobby.Unlock()

	removePlayer := func() {
		lobby.Lock()
		lobby.Players = slices.DeleteFunc(lobby.Players, func(p MultiplayerPlayer) bool { return player.Id == p.Id })
		lobby.Unlock()
		close(player.MessageChannel)
		conn.Close()
	}

	go func() {
		for {
			_, data, err := conn.ReadMessage()

			if err != nil {
				log.Println(err)
				removePlayer()
				return
			}

			var msg PlayerMessage
			err = json.Unmarshal(data, &msg)

			if err != nil {
				log.Println(err)
				removePlayer()
				return
			}

			switch msg.MsgType {
			case "renameLobby":
				if newName, ok := msg.Content.(string); ok {
					lobby.Lock()
					lobby.Name = newName
					lobby.Unlock()
				}

			}

			lobby.RLock()
			for _, ply := range lobby.Players {
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
			if len(l.Players) == 0 {
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
	logged.GET("/lobbies", getLobbies)
	router.GET("/host")
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
