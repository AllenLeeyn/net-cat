package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)
	server := server{
		logQueue:  make(chan message, 1),
		msgQueue:  make(chan message, 100),
		joinQueue: make(chan *client, 10),
		exitQueue: make(chan *client, 10),
		shutdown:  make(chan struct{})}

	switch {
	case argsLen == 0:
		server.start(":8989")
	case argsLen == 1:
		server.start(":" + args[0])
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
