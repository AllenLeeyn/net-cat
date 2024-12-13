package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
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

func runServer(portNum string) {
	listener, err := net.Listen("tcp4", portNum)
	check(err)
	defer listener.Close()

	logger("server", "Listening: "+portNum, cols["blue"])

}

func logger(userName, msg, col string) {
	now := time.Now()
	timeStamp := now.Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("[%s][%s]:%s\n", timeStamp, userName, msg)
	saveLog(msg)
	fmt.Print(col + msg)
}

func saveLog(msg string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	check(err)
	defer file.Close()

	_, err = file.WriteString(msg)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
