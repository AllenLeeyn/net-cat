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
		exitQueue: make(chan *client),
		shutdown:  make(chan struct{})}
	switch {
	case argsLen == 0:
		server.start(":8989")
	case argsLen == 1 && isPortNum(args[0]):
		server.start(args[0])
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

func isPortNum(portNum string) (isOk bool) {
	if portNum[:1] == ":" {
		if _, err := strconv.Atoi(portNum[1:]); err == nil {
			isOk = true
		}
	}
	return
}
