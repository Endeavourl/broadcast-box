package discord

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var channelId string
var discord *discordgo.Session
var streamHostUrl string

func Create() {
	token := os.Getenv("DISCORD_TOKEN")
	channelId = os.Getenv("DISCORD_CHANNEL")
	streamHostUrl = os.Getenv("STREAM_HOST_URL")

	if token == "" {
		log.Println("No DISCORD_TOKEN, Discord support disabled")
		return
	}
	if channelId == "" {
		log.Println("No DISCORD_CHANNEL, Discord support disabled")
		return
	}
	if streamHostUrl == "" {
		log.Println("No STREAM_HOST_URL, Discord support disabled")
		return
	}

	if streamHostUrl[len(streamHostUrl)-1] != '/' {
		streamHostUrl += "/"
	}

	d, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Failed to initialize Discord: " + err.Error())
	}
	discord = d
}

func SendMessage(msg string) {
	if discord == nil {
		return
	}
	_, err := discord.ChannelMessageSend(channelId, msg)
	if err != nil {
		log.Println("Failed to send Discord message: " + err.Error())
	}
}

func NotifyStream(stream string) {
	if stream[len(stream)-1] == '_' {
		// treat underscore-ending streams as hidden
		return
	}

	msgStream := "__" + strings.ReplaceAll(stream, "__", "\\__") + "__"
	urlStream := url.QueryEscape(stream)

	msg := "Stream " + msgStream + " has started. " + streamHostUrl + urlStream
	SendMessage(msg)
}

func Close() {
	if discord != nil {
		discord.Close()
	}
}
