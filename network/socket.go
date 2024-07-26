package network

import (
	"chat_server/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

//var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

type message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	When    int64  `json:"when"`
	Room    string `json:"room"`
}

type Room struct {
	Forward chan *message

	Join  chan *Client
	Leave chan *Client

	Clients map[*Client]bool

	service *service.Service
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string `json:"name"`
	Socket *websocket.Conn
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		service: service,
	}
}

func (c *Client) Read() {
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			log.Println("에러", "err", err.Error())
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
			log.Println("에러", "err", err.Error())
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
			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)

			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

const (
	SocketBufferSize  = 1024 * 4
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: messageBufferSize}

func (r *Room) ServeHTTP(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	Socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("에러", "err", err.Error())
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Println("에러", "err", err.Error())
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
