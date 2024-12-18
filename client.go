package main

import (
	"bufio"
	"net"
	"strings"
	"unicode"
)

const welcomeMsg = `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
|    ` + "`.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n" +
	"[ENTER YOUR NAME]:"

type client struct {
	conn net.Conn
	name string
}

func (s *server) addClient(cl *client) {
	isNameTaken := func(name string) bool {
		for _, cl := range s.clients {
			if cl.name == name {
				return true
			}
		}
		return false
	}
	_, err := cl.conn.Write([]byte(welcomeMsg))
	if err == nil {
		scanner := bufio.NewScanner(cl.conn)
		for scanner.Scan() {
			cl.name = strings.TrimSpace(scanner.Text())
			if cl.name == "" || !isValidName(cl.name) {
				cl.conn.Write([]byte("Invalid entry. Try again: "))
				continue
			}
			if isNameTaken(cl.name) {
				cl.conn.Write([]byte("Name taken. Try again: "))
				continue
			}
			if len(s.clients) >= maxConn {
				cl.conn.Write([]byte("Server full. Try again later.\n"))
				break
			}
			s.joinQueue <- cl
			s.msgQueue <- message{from: "server",
				body: []byte(cl.name + " has joined the chat.\n")}
			return
		}
	}
}

func (s *server) removeClient(cl *client) {

	for i, c := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			s.msgQueue <- message{from: "server",
				body: []byte(cl.name + " has leaved the chat.\n")}
		}
	}
	cl.conn.Close()
}

func isValidName(name string) bool {
	for _, rn := range name {
		if !unicode.IsPrint(rn) {
			return false
		}
	}
	return true
}
