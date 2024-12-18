package main

import (
	"fmt"
	"log"
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

func formatMsg(userName, msg, col string) []byte {
	now := time.Now()
	timeStamp := now.Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s[%s][%s]:%s%s", col, timeStamp, userName, msg, cols["reset"])
	return []byte(msg)
}

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
