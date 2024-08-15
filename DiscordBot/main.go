package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

const TOKEN_ENV = "DISCORD_BOT_TOKEN"
const APP_ENV = "DISCORD_BOT_APP"
const GUILD_ENV = "DISCORD_BOT_GUILD"
const SERVER_URL_ENV = "WEB_SERVER_URL"
const SERVER_TOKEN_ENV = "WEB_SERVER_TOKEN"

var (
	Token       = os.Getenv(TOKEN_ENV)
	App         = os.Getenv(APP_ENV)
	Guild       = os.Getenv(GUILD_ENV)
	ServerUrl   = os.Getenv(SERVER_URL_ENV)
	ServerToken = os.Getenv(SERVER_TOKEN_ENV)
)

func handleSimplify(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionData) {
	option := data.Options[0]
	if option.Name != "image" || !option.BoolValue() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Missing image!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	attachmentId := option.Value.(string)
	attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentId]
	attachmentUrl := attachment.URL

	errorRespond := func(logString string, responseString string) {
		log.Println(logString)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: responseString,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	res, err := http.Get(attachmentUrl)
	if err != nil {
		errorRespond(fmt.Sprintf("Couldn't get image at '%s'", attachmentUrl), "Couldn't get the image you attached.")
		return
	}

	var buff bytes.Buffer
	w := multipart.NewWriter(&buff)
	fw, err := w.CreateFormFile("image", attachment.Filename)
	if err != nil {
		errorRespond(fmt.Sprintf("Couldn't write image at '%s' to form.", attachmentUrl), "Couldn't send the image to the server.")
		return
	}

	if _, err = io.Copy(fw, res.Body); err != nil {
		errorRespond(fmt.Sprintf("Couldn't write image at '%s' to form.", attachmentUrl), "Couldn't send the image to the server.")
		return
	}
	res.Body.Close()

	w.Close()

	req, err := http.NewRequest("POST", ServerUrl, &buff)
	if err != nil {
		errorRespond(fmt.Sprintf("Couldn't make request to server with image at '%s'", attachmentUrl), "Couldn't send the image to the server.")
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		errorRespond(fmt.Sprintf("Request to server containing image at '%s' has failed", attachmentUrl), "Request to server has failed.")
		return
	}

	if res.StatusCode != http.StatusOK {
		if b, err := io.ReadAll(res.Body); err == nil {
			if json.Valid(b) {
				m := make(map[string]any)
				if err = json.Unmarshal(b, &m); err == nil {
					if errorMessage, ok := m["errorMessage"]; ok { // THE GREAT PYRAMID
						errorRespond(fmt.Sprintf("Request to server containing image at '%s' has errored with code %d: %v", attachmentUrl, res.StatusCode, errorMessage), fmt.Sprintf("Error: %v", errorMessage))
						return
					}
				}
			}
		}
		errorRespond(fmt.Sprintf("Request to server containing image at '%s' has errored with code %d", attachmentUrl, res.StatusCode), "The server returned an error.")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("File %s was successfully uploaded and simplified.", attachment.Filename),
		},
	})
}

func openServerWebsocket(session *discordgo.Session) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ServerUrl, Path: "/api/ws"}
	q := u.Query()
	q.Add("token", ServerToken)
	u.RawQuery = q.Encode()
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Couldn't open a connection to the main server: %s", err)
	}
	log.Println("Connected to main server!")

	go func() {
		for {
			msg := make(map[string]any)
			err := conn.ReadJSON(&msg)
			if err != nil {
				e, ok := err.(*websocket.CloseError)
				if !ok || e.Code != websocket.CloseNormalClosure {
					log.Println(err)
				}
				break
			}
			log.Println(msg)
		}
		conn.Close()
	}()

	return conn
}

func main() {
	if Token == "" {
		log.Fatalf("Bot token not set in environment var '%s'", TOKEN_ENV)
	}
	if App == "" {
		log.Fatalf("Application id not set in environment var '%s'", APP_ENV)
	}
	if Guild == "" {
		log.Fatalf("Guild id not set in environment var '%s'", GUILD_ENV)
	}
	if ServerUrl == "" {
		log.Fatalf("Server URL not set in environment var '%s'", SERVER_URL_ENV)
	}

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		if data.Name != "simplify" {
			return
		}

		handleSimplify(s, i, &data)
	})

	// TODO: add simplify slash command creation
	// TODO: relay simplify requests to bot channel or just give results of post request idk

	err = session.Open()
	if err != nil {
		log.Fatalf("Couldn't open a connection to Discord: %s", err)
	}
	log.Println("Connected to Discord!")
	defer session.Close()

	serverConn := openServerWebsocket(session)
	defer serverConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigch
	log.Println("Shutting down...")
}
