package main

import (
	"fmt"
	"os"
	"strconv"
)

var cols = map[string]string{
	"black":   "\033[38;2;000;000;000m",
	"red":     "\033[38;2;255;000;000m",
	"green":   "\033[38;2;000;255;000m",
	"yellow":  "\033[38;2;255;255;000m",
	"blue":    "\033[38;2;000;000;255m",
	"magenta": "\033[38;2;255;000;255m",
	"cyan":    "\033[38;2;000;255;255m",
	"white":   "\033[38;2;255;255;255m",
	"orange":  "\033[38;2;255;165;000m",
	"reset":   "\033[00m",
}

func main() {
	args := os.Args[1:]
	argsLen := len(args)
	server := server{msgQueue: make(chan []byte),
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
