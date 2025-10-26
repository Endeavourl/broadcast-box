package discord

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/glimesh/broadcast-box/internal/util"
	"github.com/glimesh/broadcast-box/internal/webrtc"
)

type streamInfo struct {
	active    bool
	clients   int
	lastStart time.Time
	lastEnd   time.Time
}

var (
	channelId     string
	streamHostUrl string
	discord       *discordgo.Session

	streamsLock sync.Mutex
	streams     map[string]*streamInfo
)

func Create() {
	streams = map[string]*streamInfo{}

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
		fmt.Println("Failed to initialize Discord:", err)
		return
	}
	discord = d

	webrtc.AddStreamListener(streamListener)
	webrtc.AddStreamClientListener(streamClientListener)

	err = d.Open()
	if err != nil {
		fmt.Println("Error opening Discord connection:", err)
		return
	}
	setOfflineStatus()
}

func streamListener(streamName string, started bool) {
	streamsLock.Lock()
	defer streamsLock.Unlock()

	si := getStreamInfo(streamName)
	si.active = started
	si.clients = 0
	if started {
		si.lastStart = time.Now()
		notifyStream(streamName, si)
	} else {
		si.lastEnd = time.Now()
	}
	updateStatus()
}

func streamClientListener(streamName string, clients int) {
	streamsLock.Lock()
	defer streamsLock.Unlock()

	getStreamInfo(streamName).clients = clients
	updateStatus()
}

func getStreamInfo(streamName string) *streamInfo {
	si, ok := streams[streamName]
	if !ok {
		si = &streamInfo{}
		streams[streamName] = si
	}
	return si
}

func updateStatus() {
	status := ""
	for stream, si := range streams {
		if !si.active || util.IsHiddenStream(stream) {
			continue
		}
		status += stream + ": " + strconv.Itoa(si.clients) + ", "
	}
	if len(status) > 0 {
		status = status[:len(status)-2]
		go setStreamStatus(status)
	} else {
		go setOfflineStatus()
	}
}

func notifyStream(stream string, streamInfo *streamInfo) {
	if util.IsHiddenStream(stream) {
		return
	}

	if time.Since(streamInfo.lastEnd).Seconds() <= 60 {
		// do not notify about restarts (previous stream ended < minute ago)
		return
	}

	msgStream := "__" + strings.ReplaceAll(stream, "__", "\\__") + "__"
	urlStream := url.QueryEscape(stream)

	msg := "Stream " + msgStream + " has started. " + streamHostUrl + urlStream
	go sendMessage(msg)
}

func sendMessage(msg string) {
	if discord == nil {
		return
	}
	_, err := discord.ChannelMessageSend(channelId, msg)
	if err != nil {
		fmt.Println("Failed to send Discord message:", err)
	}
}

func setStreamStatus(status string) {
	if discord == nil {
		return
	}
	err := discord.UpdateWatchStatus(0, status)
	if err != nil {
		fmt.Println("Failed to set Discord status:", err)
	}
}

func setOfflineStatus() {
	if discord == nil {
		return
	}
	err := discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "invisible",
	})
	if err != nil {
		fmt.Println("Failed to set Discord status to offline:", err)
	}
}

func Close() {
	if discord != nil {
		setOfflineStatus()
		discord.Close()
	}
}
