package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const maxConn = 1

type server struct {
	clients   []*client
	log       *os.File
	logQueue  chan message
	history   [][]byte
	msgQueue  chan message
	joinQueue chan *client
	exitQueue chan *client
}

func (s *server) start(portNum string) {
	server, err := net.Listen("tcp", portNum)
	check(err)
	defer server.Close()

	fileName := fmt.Sprintf("sys@%s-[%s].log", portNum, getTimeStamp())
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0o644)
	check(err)
	s.log = file
	defer file.Close()

	s.logQueue <- message{from: "server", body: []byte("Listening on port " + portNum)}

	go s.listener()
	go s.broadcaster()

	for {
		conn, err := server.Accept()
		if err != nil {
			conn.Close()
		}
		go s.handlerConnection(conn)

	}
}

func (s *server) handlerConnection(conn net.Conn) {
	s.logQueue <- message{from: "server",
		body: []byte("connecting " + conn.RemoteAddr().String())}

	cl := &client{conn: conn}
	s.addClient(cl)

	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		if isValidEntry(scanner.Text()) {
			s.msgQueue <- message{from: cl.name, body: []byte(scanner.Text())}
		}
	}
	s.exitQueue <- cl
}

func (s *server) listener() {
	for {
		select {
		case msg := <-s.logQueue:
			msgPretty := formatMsg(getTimeStamp(), msg, "")
			_, err := s.log.Write(msgPretty)
			check(err)
			fmt.Print(string(msgPretty))

		case cl := <-s.joinQueue:
			if len(s.clients) >= maxConn {
				s.logQueue <- message{from: "server",
					body: []byte("Server full. Unable to connect " + cl.conn.RemoteAddr().String())}
				cl.conn.Write([]byte("Server full. Try again later.\n"))
				cl.conn.Close()
				continue
			}
			s.clients = append(s.clients, cl)
			s.msgQueue <- message{from: "server",
				body: []byte(cl.name + " has joined the chat.")}

		case cl := <-s.exitQueue:
			s.removeClient(cl)
		}
	}
}

func (s *server) broadcaster() {
	for {
		msg := <-s.msgQueue
		timeStamp := getTimeStamp()
		color := getMsgColor(msg)
		msgPretty := formatMsg(timeStamp, msg, color)

		s.logQueue <- msg
		s.history = append(s.history, msgPretty)

		for _, cl := range s.clients {
			if msg.from == cl.name {
				msgPretty = formatMsg(timeStamp, msg, "green")
			}
			_, err := cl.conn.Write(msgPretty)
			if err != nil {
				return
			}
			if msg.from == cl.name {
				msgPretty = formatMsg(timeStamp, msg, color)
			}
		}
	}
}
