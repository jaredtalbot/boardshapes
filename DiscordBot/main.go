package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

const TOKEN_ENV = "DISCORD_BOT_TOKEN"
const APP_ENV = "DISCORD_BOT_APP"
const GUILD_ENV = "DISCORD_BOT_GUILD"
const CHANNEL_ENV = "DISCORD_BOT_CHANNEL"
const SERVER_URL_ENV = "WEB_SERVER_URL"
const SERVER_TOKEN_ENV = "WEB_SERVER_TOKEN"

var (
	Token       = os.Getenv(TOKEN_ENV)
	App         = os.Getenv(APP_ENV)
	Channel     = os.Getenv(CHANNEL_ENV)
	Guild       = os.Getenv(GUILD_ENV)
	ServerUrl   = os.Getenv(SERVER_URL_ENV)
	ServerToken = os.Getenv(SERVER_TOKEN_ENV)
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func handleSimplify(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionData) {
	option := data.Options[0]
	if option.Name != "image" || option.Type != discordgo.ApplicationCommandOptionAttachment {
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

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"image", quoteEscaper.Replace(attachment.Filename)))
	h.Set("Content-Type", "image/png")
	fw, err := w.CreatePart(h)
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

	req, err := http.NewRequest("POST", ServerUrl+"/api/simplify", &buff)
	if err != nil {
		errorRespond(fmt.Sprintf("Couldn't make request to server with image at '%s'", attachmentUrl), "Couldn't send the image to the server.")
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
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
						log.Printf("Request to server containing image at '%s' has errored with code %d: %v", attachmentUrl, res.StatusCode, errorMessage)
						s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
							Content: fmt.Sprintf("Error: %v", errorMessage),
						})
						return
					}
				}
			}
		}
		log.Printf("Request to server containing image at '%s' has errored with code %d", attachmentUrl, res.StatusCode)
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: fmt.Sprintf("There was an error with uploading %s to the server.", attachment.Filename),
		})
		return
	}
	s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("File %s was successfully uploaded and simplified.", attachment.Filename),
	})
}

type WebServerAttachedFile struct {
	Name          string         `json:"name"`
	ContentType   string         `json:"contentType"`
	Base64Content string         `json:"base64Content"`
	Meta          map[string]any `json:"meta,omitempty"`
}

type WebServerMessage struct {
	Type        string                  `json:"type"`
	Content     string                  `json:"content,omitempty"`
	Attachments []WebServerAttachedFile `json:"attachments,omitempty"`
	Timestamp   string                  `json:"timestamp,omitempty"`
}

var serverConn *websocket.Conn

var retryDelays = []int{5, 10, 30, 60}

func handleWebServerMessage(session *discordgo.Session, msg *WebServerMessage) {
	var discordMsgFiles []*discordgo.File
	var discordMsgEmbeds []*discordgo.MessageEmbed
	if msg.Attachments != nil {
		discordMsgFiles = make([]*discordgo.File, len(msg.Attachments))
		for i, v := range msg.Attachments {
			discordMsgFiles[i] = &discordgo.File{
				Name:        v.Name,
				ContentType: v.ContentType,
				Reader:      base64.NewDecoder(base64.StdEncoding, strings.NewReader(v.Base64Content)),
			}
		}
		if msg.Type == "simplify" {
			discordMsgEmbeds = make([]*discordgo.MessageEmbed, len(msg.Attachments))
			for i, v := range msg.Attachments {
				var footer strings.Builder
				footer.WriteString("Region Count: ")
				if regionCount, ok := v.Meta["regionCount"]; ok {
					footer.WriteString(fmt.Sprint(regionCount))
				}
				discordMsgEmbeds[i] = &discordgo.MessageEmbed{
					Title:       "Simplified Image",
					Description: v.Name,
					Image: &discordgo.MessageEmbedImage{
						URL: "attachment://" + v.Name,
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: footer.String(),
					},
					Timestamp: msg.Timestamp,
					Color:     0x237feb,
				}
			}
		}
	}

	discordMsg := &discordgo.MessageSend{Content: msg.Content, Files: discordMsgFiles, Embeds: discordMsgEmbeds}
	_, err := session.ChannelMessageSendComplex(Channel, discordMsg)
	if err != nil {
		log.Printf("Couldn't send message from main server [%s]", err)
	}
}

// TODO: better name?
func handleServerEvents(session *discordgo.Session, u *url.URL) {
	retries := 0
	suppressRetryMessages := false
	for {
		// try to connect
		serverConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			retryDelay := retryDelays[min(retries, len(retryDelays)-1)]
			retries++
			if !suppressRetryMessages {
				log.Printf("Couldn't open a connection to the main server [%s]", err)
				log.Printf("Retrying in %d seconds", retryDelay)
			}
			time.Sleep(time.Second * time.Duration(retryDelay))
			continue
		}
		// connected
		log.Println("Connected to main server!")
		retries = 0
		suppressRetryMessages = false
	handleMessages:
		for {
			msg := new(WebServerMessage)
			err := serverConn.ReadJSON(msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Connection with main server was closed normally.")
					suppressRetryMessages = true
					break handleMessages
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Disconnected from main server [%s]", err)
					break handleMessages
				}
				switch err.(type) {
				case *json.SyntaxError, *json.UnmarshalTypeError:
					log.Printf("Bad data received from main server [%s]", err)
				default:
					log.Println(err)
					break handleMessages
				}
			}
			go handleWebServerMessage(session, msg)
		}
	}
}

func openServerWebsocket(session *discordgo.Session) error {
	u, err := url.Parse(ServerUrl)
	if err != nil {
		log.Printf("Couldn't parse server URL [%s]", err)
		return err
	}
	u.Path = "/api/ws"
	u.Scheme = "ws"
	if strings.Contains(ServerUrl, "https") {
		u.Scheme = "wss"
	}
	q := u.Query()
	q.Add("token", ServerToken)
	u.RawQuery = q.Encode()

	go handleServerEvents(session, u)

	return nil
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "simplify",
		Description: "Try out simplifying an image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "image",
				Description: "Contents of the message",
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Required:    true,
			},
		},
		Type: discordgo.ChatApplicationCommand,
	},
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
	if Channel == "" {
		log.Fatalf("Channel id not set in environment var '%s'", CHANNEL_ENV)
	}
	if ServerUrl == "" {
		log.Fatalf("Server URL not set in environment var '%s'", SERVER_URL_ENV)
	}
	if ServerToken == "" {
		log.Fatalf("Server Token not set in environment var '%s'", SERVER_TOKEN_ENV)
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

	_, err = session.ApplicationCommandBulkOverwrite(App, Guild, commands)
	if err != nil {
		log.Printf("Could not register commands: %s", err)
		return
	}

	err = session.Open()
	if err != nil {
		log.Printf("Couldn't open a connection to Discord: %s", err)
		return
	}
	log.Println("Connected to Discord!")
	defer session.Close()

	err = openServerWebsocket(session)
	if err != nil {
		return
	}
	defer func() {
		if serverConn != nil {
			defer serverConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		}
	}()

	defer log.Println("Shutting down...")
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigch
}
