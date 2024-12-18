package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)
	server := server{msgQueue: make(chan message, 10),
		joinQueue: make(chan *client, 10),
		exitQueue: make(chan *client, 10),
		shutdown:  make(chan struct{}, 1)}

	switch {
	case argsLen == 0:
		server.start(":8989")
	case argsLen == 1 && isPortNum(args[0]):
		server.start(":" + args[0])
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

func isPortNum(portNum string) (isOk bool) {
	if _, err := strconv.Atoi(portNum[1:]); err == nil {
		isOk = true
	}
	return
}
