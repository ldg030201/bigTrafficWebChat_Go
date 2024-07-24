package main

import (
	"chat_server/config"
	"chat_server/repository"
	"flag"
	"fmt"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()

	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		fmt.Println(rep)
	}

	/*n := network.NewServer()
	n.StartServer()*/
}
