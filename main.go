package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	switch {
	case argsLen == 0:
		runServer(":8989")
	case argsLen == 1 && isPortNum(args[0]):
		runServer(args[0])
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
