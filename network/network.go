package network

import (
	"chat_server/service"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	engin *gin.Engine

	service *service.Service

	port string
	ip   string
}

func NewServer(service *service.Service, port string) *Server {
	s := &Server{
		engin:   gin.New(),
		service: service,
		port:    port,
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

	registerServer(s)

	return s
}

func (s *Server) StartServer() error {
	s.setServerInfo()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT)

	go func() {
		<-channel

		if err := s.service.ServerSet(s.ip+s.port, false); err != nil {
			log.Printf("server error")
		}

		type ServerInfoEvent struct {
			IP     string
			Status bool
		}

		e := &ServerInfoEvent{IP: s.ip + s.port, Status: false}
		ch := make(chan kafka.Event)

		if v, err := json.Marshal(e); err != nil {
			log.Printf("json marshal error")
		} else if result, err := s.service.PublishEvent("chat", v, ch); err != nil {
			log.Println("end event kafka error")
		} else {
			log.Println("success event", result)
		}

		os.Exit(1)
	}()

	log.Println("서버 시작")
	return s.engin.Run(s.port)
}

func (s *Server) setServerInfo() {
	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(err.Error())
	} else {
		var ip net.IP

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					break
				}
			}
		}

		if ip == nil {
			panic("no ip")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err.Error())
			} else {
				s.ip = ip.String()
			}
		}
	}
}
