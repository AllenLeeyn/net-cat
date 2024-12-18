package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

const maxConn = 2

type server struct {
	clients   []*client
	msgQueue  chan message
	log       *os.File
	history   [][]byte
	joinQueue chan *client
	exitQueue chan *client
	shutdown  chan struct{}
}

type message struct {
	from string
	body []byte
}

func (s *server) start(portNum string) {
	server, err := net.Listen("tcp4", portNum)
	check(err)

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
	check(err)
	s.log = file
	defer file.Close()

	s.msgQueue <- message{from: "server", body: []byte("Listening on port " + portNum + "\n")}
	go s.listener()
	go s.broadcaster()

	for {
		conn, err := server.Accept()
		if err != nil {
			conn.Close()
			continue
		}
		go s.handlerConnection(conn)
	}

	//<-s.shutdown
}

func (s *server) handlerConnection(conn net.Conn) {
	cl := &client{conn: conn,
		from: make(chan []byte, 10),
		exit: make(chan struct{}, 1)}

	s.addClient(cl)

	<-cl.exit
	s.exitQueue <- cl
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
			go cl.getFrom()
			return
		}
	}
	cl.exit <- struct{}{}
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
	close(cl.from)
	close(cl.exit)
}

func (s *server) listener() {
	for {
		select {
		case cl := <-s.joinQueue:
			s.clients = append(s.clients, cl)
			s.getHistory(cl)
		case cl := <-s.exitQueue:
			s.removeClient(cl)
		default:
			for _, cl := range s.clients {
				select {
				case msg := <-cl.from:
					fromMsg := message{from: cl.name, body: msg}
					s.msgQueue <- fromMsg
				default:
				}
			}
		}
	}
}

func (s *server) getHistory(cl *client) {
	for _, msg := range s.history {
		_, err := cl.conn.Write(msg)
		if err != nil {
			cl.exit <- struct{}{}
			break
		}
	}
}

func (s *server) broadcaster() {
	for {
		msg := <-s.msgQueue
		col := cols["yellow"]
		if msg.from != "server" {
			col = cols["blue"]
		}
		msgPretty := formatMsg(msg.from, string(msg.body), col)

		_, err := s.log.Write(msgPretty)
		check(err)
		s.history = append(s.history, msgPretty)

		fmt.Print(string(msgPretty))

		for _, cl := range s.clients {
			if msg.from == cl.name {
				msgPretty = formatMsg(msg.from, string(msg.body), cols["green"])
			}
			_, err := cl.conn.Write(msgPretty)
			if err != nil {
				cl.exit <- struct{}{}
				return
			}
			if msg.from == cl.name {
				msgPretty = formatMsg(msg.from, string(msg.body), col)
			}
		}
	}
}

func isValidName(name string) bool {
	for _, rn := range name {
		if !unicode.IsPrint(rn) {
			return false
		}
	}
	return true
}
