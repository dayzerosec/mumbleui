package socketserver

import (
	io "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"log"
	"net/http"
)

type WSS struct {
	io *io.Server
	OnMessage func(msg SocketMessage) interface{}
	StaticDir string
	HostAddr string
}

type SocketMessage struct {
	Action string `json:"action"`
	Id string `json:"id"`
	Data interface{} `json:"data"`
}


func (s *WSS) Broadcast(method string, args interface{}) {
	if s.io == nil {
		return
	}
	s.io.BroadcastToAll(method, args)
}

func (s *WSS) Init() {

}
func (s *WSS) onEvent(c *io.Channel,  msg SocketMessage) {
	response := SocketMessage{
		Action: msg.Action,
		Id:     msg.Id,
		Data:   s.OnMessage(msg),
	}
	if err := c.Emit("action", response); err != nil {
		log.Println("[!]", err)
	}
}

func (s *WSS) Start() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(s.StaticDir)))

	s.io = io.NewServer(transport.GetDefaultWebsocketTransport())
	serveMux.Handle("/socket.io/", s.io)
	if err := s.io.On("action", s.onEvent); err != nil {
		panic(err)
	}

	log.Printf("Serving on " + s.HostAddr)
	if err := http.ListenAndServe(s.HostAddr, serveMux); err != nil {
		panic(err)
	}
}