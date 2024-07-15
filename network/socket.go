package network

import (
	"chat_server/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Room struct {
	Forward chan *message

	Join  chan *Client
	Leave chan *Client

	Clients map[*Client]bool
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) RunInit() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client

	defer func() { r.Leave <- client }()
}
