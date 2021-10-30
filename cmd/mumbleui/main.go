package main

import (
	"crypto/tls"
	"flag"
	"layeh.com/gumble/gumble"
	_ "layeh.com/gumble/opus"
	"log"
	"mumbleui/pkg/mumbletracker"
	"mumbleui/pkg/socketserver"
	"net"
	"os"
	"time"
)

type SpeakerEvent struct {
	Name string `json:"name"`
	Speaking bool `json:"is_speaking"`
}

var server = socketserver.WSS{
	OnMessage: OnSocketMessage,
	HostAddr: "0.0.0.0:8080",
	StaticDir: "./web",
}

var audioListener = &mumbletracker.AudioListener{
	Frequency:           100 * time.Millisecond,
	OnStartSpeaking:     OnStartSpeaking,
	OnStopSpeaking:      OnStopSpeaking,
}
var eventListener = mumbletracker.NewEventListener(OnJoin, OnLeave)
var currentUserList []string

func OnSocketMessage(message socketserver.SocketMessage) interface{} {
	switch message.Action {
	case "user-list":
		return currentUserList
	}
	return nil
}

func OnJoin(user *gumble.User, list []string) {
	server.Broadcast("join", user.Name)
	currentUserList = list
	log.Println(user.Name, "joined")
}

func OnLeave(user *gumble.User, list []string) {
	log.Printf("[%s] Left.", user.Name)
	server.Broadcast("leave", user.Name)
	currentUserList = list
	log.Println(user.Name, "disconnected")
}

func OnStartSpeaking(user *gumble.User) {
	server.Broadcast("broadcast", SpeakerEvent{
		Name:     user.Name,
		Speaking: true,
	})
}

func OnStopSpeaking(user *gumble.User) {
	server.Broadcast("broadcast", SpeakerEvent{
		Name:     user.Name,
		Speaking: false,
	})
}

func getArguments(addr, username, password *string, verify *bool) {
	flag.StringVar(addr, "addr", "", "server address (host:port)")
	flag.StringVar(password, "pw", "", "server password")
	flag.StringVar(username, "user", "", "username")
	flag.BoolVar(verify, "verify", true, "Whether or not to verify the TLS certificate")
	flag.Parse()

	if *addr == "" {
		if val, found := os.LookupEnv("MUMBLE_ADDR"); found {
			*addr = val
		}
	}
	if *username == "" {
		if val, found := os.LookupEnv("MUMBLE_USER"); found {
			*username = val
		}
	}
	if *password == "" {
		if val, found := os.LookupEnv("MUMBLE_PW"); found {
			*password = val
		}
	}
}

func main() {
	var addr, username, password string
	var verify bool
	getArguments(&addr, &username, &password, &verify)

	config := gumble.NewConfig()
	config.Username = username
	config.Password = password

	config.Attach(eventListener)

	audioListener.Frequency = 75 * time.Millisecond
	config.AttachAudio(audioListener)

	tlscfg := &tls.Config{InsecureSkipVerify: !verify}
	_, err := gumble.DialWithDialer(new(net.Dialer), addr, config, tlscfg)
	if err != nil {
		panic(err)
	}
	server.Start()
}

