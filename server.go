package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const maxConn = 10

type server struct {
	self      *client
	clients   []*client
	log       *os.File
	logQueue  chan message
	history   [][]byte
	msgQueue  chan message
	joinQueue chan *client
	exitQueue chan *client
}

// s.start() starts listening for request at portNum,
// openFile() for logging and s.Accept() in coming connections.
// It starts the s.listener and s.broadcaster as go routines.
// For each conn s.Accpet(), a go s.handleConnection routine is used to handle it.
func (s *server) start(portNum string) {
	server, err := net.Listen("tcp", portNum)
	check(err)
	defer server.Close()

	fileName := fmt.Sprintf("%v@%v.log",
		time.Now().Format("20060102_150405"), portNum[1:])
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0o644)
	check(err)
	s.log = file
	defer file.Close()

	s.self = &client{name: "server", color: colors[time.Now().Second()%12]}
	s.logQueue <- message{from: s.self, body: []byte("Listening on port " + portNum)}

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

// s.handleConnection() tries to processClient() and Scan() for incoming messages.
// If processClient() fails, the conn will be close and an error will occur here.
// Incoming messages are written to the s.msgQueue.
func (s *server) handlerConnection(conn net.Conn) {
	s.logQueue <- message{from: s.self,
		body: []byte("connecting " + conn.RemoteAddr().String())}

	cl := &client{conn: conn, color: colors[time.Now().Second()%16]}
	s.processClient(cl)

	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		if isValidEntry(scanner.Text()) {
			s.msgQueue <- message{from: cl, body: []byte(scanner.Text())}
		}
	}
	s.exitQueue <- cl
}

// s.listenser() listens to acitivity on s.logQueue, s.joinQueue and s.exitQueue and handles them.
func (s *server) listener() {
	for {
		select {
		case msg := <-s.logQueue:
			msgPretty := formatMsg(msg)
			_, err := s.log.Write(msgPretty)
			check(err)
			fmt.Print(string(msgPretty))

			// registers connecting client if server is not full.
			// Prints history to connecting client and
			// notify other clients of new client.
		case cl := <-s.joinQueue:
			if len(s.clients) >= maxConn {
				s.logQueue <- message{from: s.self,
					body: []byte("Server full. Unable to connect " + cl.conn.RemoteAddr().String())}
				cl.conn.Write([]byte("Server full. Try again later.\n"))
				cl.conn.Close()
				continue
			}
			s.clients = append(s.clients, cl)
			if err := s.printHistory(cl); err != nil {
				cl.conn.Close()
				continue
			}
			s.msgQueue <- message{from: s.self,
				body: []byte(cl.name + " has joined the chat.")}

		case cl := <-s.exitQueue:
			s.removeClient(cl)
		}
	}
}

// s.broadcaster() grabs msg from msgQueue (if any).
// The msg will be logged, save to history, and send to all client.
func (s *server) broadcaster() {
	for {
		msg := <-s.msgQueue
		msgPretty := formatMsg(msg)

		s.logQueue <- msg
		s.history = append(s.history, msgPretty)

		for _, cl := range s.clients {
			_, err := cl.conn.Write(msgPretty)
			if err != nil {
				cl.conn.Close()
			}
		}

	}
}
