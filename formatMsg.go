package main

import (
	"fmt"
	"time"
)

var colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[00m",
}

type message struct {
	from string
	body []byte
}

func getTimeStamp() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func getMsgColor(msg message) string {
	color := "yellow"
	if msg.from != "server" {
		color = "blue"
	}
	return color
}

func formatMsg(timeStamp string, msg message, color string) []byte {
	msgPretty := fmt.Sprintf("%s[%s][%s]:%s%s\n",
		colors[color], timeStamp, msg.from, msg.body, colors["reset"])

	return []byte(msgPretty)
}
