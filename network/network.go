package network

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	n := &Network{
		engin: gin.New(),
	}

	return n
}

func (n *Network) StartServer() error {
	log.Println("서버 시작")
	return n.engin.Run(":8080")
}
