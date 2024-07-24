package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

//var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

type message struct {
	Name    string
	Message string
	When    int64
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

func (c *Client) Read() {
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			panic(err)
		} else {
			msg.When = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()

	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) Run() {
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

const (
	SocketBufferSize  = 1024 * 16
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: messageBufferSize}

func (r *Room) ServeHTTP(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	Socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Socket: Socket,
		Send:   make(chan *message, messageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- client

	defer func() { r.Leave <- client }()

	go client.Write()

	client.Read()
}
