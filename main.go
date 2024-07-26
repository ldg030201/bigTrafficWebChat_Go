package main

import (
	"chat_server/config"
	"chat_server/network"
	"chat_server/repository"
	"chat_server/service"
	"flag"
	"log"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()

	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		log.Println("에러", "err", err.Error())
	} else {
		n := network.NewServer(service.NewService(rep), *port)
		n.StartServer()
	}
}
