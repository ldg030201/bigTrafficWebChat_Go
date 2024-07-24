package network

import (
	"chat_server/repository"
	"chat_server/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	engin *gin.Engine
}

func NewServer(service *service.Service, repository *repository.Repository, port string) *Server {
	s := &Server{
		engin: gin.New(),
	}

	s.engin.Use(gin.Logger())
	s.engin.Use(gin.Recovery())
	s.engin.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	registerServer(s.engin)

	return s
}

func (n *Server) StartServer() error {
	log.Println("서버 시작")
	return n.engin.Run(":8080")
}
