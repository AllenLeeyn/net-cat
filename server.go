package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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
	cl := &client{conn: conn}

	s.addClient(cl)
	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		msg := message{from: cl.name, body: []byte(scanner.Text() + "\n")}
		s.msgQueue <- msg
	}
	s.exitQueue <- cl
}

func (s *server) listener() {
	for {
		select {
		case cl := <-s.joinQueue:
			s.clients = append(s.clients, cl)
			s.getHistory(cl)
		case cl := <-s.exitQueue:
			s.removeClient(cl)
		}
	}
}

func (s *server) getHistory(cl *client) {
	for _, msg := range s.history {
		_, err := cl.conn.Write(msg)
		if err != nil {
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
				return
			}
			if msg.from == cl.name {
				msgPretty = formatMsg(msg.from, string(msg.body), col)
			}
		}
	}
}
