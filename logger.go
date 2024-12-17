package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

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
