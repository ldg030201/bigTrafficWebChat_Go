package main

import "chat_server/network"

func main() {
	n := network.NewServer()
	n.StartServer()
}
