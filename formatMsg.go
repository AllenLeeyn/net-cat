package main

import (
	"fmt"
	"time"
)

const resetColor = "\033[00m"

var colors = []string{
	"\033[31m", //red
	"\033[32m", //green
	"\033[33m", //yellow
	"\033[34m", //blue
	"\033[35m", //magenta
	"\033[36m", //cyan
	"\033[91m", //high intense red
	"\033[92m", //high intense green
	"\033[93m", //high intense yellow
	"\033[94m", //high intense blue
	"\033[95m", //high intense magenta
	"\033[96m", //high intense cyan
}

type message struct {
	from *client
	body []byte
}

// formatMsg() format messages with timestamp and sender
func formatMsg(msg message) []byte {
	now := time.Now()
	timeStamp := now.Format("2006-01-02 15:04:05")
	msgPretty := fmt.Sprintf("%s[%s][%s]:%s%s\n",
		msg.from.color, timeStamp, msg.from.name, msg.body, resetColor)

	return []byte(msgPretty)
}

// isValidEntry() checks if entry is printable ascii and åäö
func isValidEntry(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	for _, rn := range entry {
		if !((rn >= 32 && rn <= 126) ||
			rn == 'å' || rn == 'ä' || rn == 'ö' ||
			rn == 'Å' || rn == 'Ä' || rn == 'Ö') {
			return false
		}
	}
	return true
}
