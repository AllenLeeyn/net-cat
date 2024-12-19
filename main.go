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
		logQueue:  make(chan message, 100),
		msgQueue:  make(chan message, 100),
		joinQueue: make(chan *client, 10),
		exitQueue: make(chan *client, 10)}

	switch {
	case argsLen == 0:
		server.start(":8989")
	case argsLen == 1:
		server.start(":" + args[0])
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
	fmt.Print("here")
}

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
