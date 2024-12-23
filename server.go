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
	shutdown  chan struct{}
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

	s.self = &client{name: "server", color: "\033[7m" + colors[time.Now().Second()%12]}
	s.logQueue <- message{from: s.self, body: []byte("Listening on port " + portNum)}

	go s.listener()
	go s.broadcaster()

	// make server.Accept() as a go rountine.
	// Here if an error occurs (regardless of client or server side),
	// the program just exits the go routine
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println("Exiting server connector...")
				break
			}
			go s.handlerConnection(conn)
		}
		fmt.Print()
	}()

	// implement shutdown listener
	fmt.Println("Enter 'quit' to shutdown")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() != nil || scanner.Text() == "quit" {
			break
		}
	}
	s.stop()
}

// s.stop() annouces shutdown and close all connection and exit go routines
func (s *server) stop() {
	s.msgQueue <- message{from: s.self,
		body: []byte("Server shutting down...")}
	time.Sleep(1 * time.Second) // wait for message to be send
	for _, cl := range s.clients {
		cl.conn.Close()
	}
	s.msgQueue <- message{body: []byte("\033exit")}
	s.shutdown <- struct{}{}
	s.log.Close()
}

// s.handleConnection() tries to processClient() and Scan() for incoming messages.
// If processClient() fails, the conn will be close and an error will occur here.
// Incoming messages are written to the s.msgQueue.
func (s *server) handlerConnection(conn net.Conn) {
	s.logQueue <- message{from: s.self,
		body: []byte("connecting " + conn.RemoteAddr().String())}

	cl := &client{conn: conn, color: colors[time.Now().Second()%12]}
	s.processClient(cl)

	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		if !isValidEntry(scanner.Text()) {
			continue
		}
		if scanner.Text() == "--rename" {
			if !s.setClientName(cl) {
				break
			}
			continue
		}
		if scanner.Text() == "--recolor" {
			oldColor := cl.color
			cl.color = colors[time.Now().Second()%12]
			s.msgQueue <- message{from: cl,
				body: []byte(cl.name + " changed color from " + oldColor + "this " + cl.color + "to this")}
			continue
		}
		s.msgQueue <- message{from: cl, body: []byte(scanner.Text())}
	}
	s.exitQueue <- cl
}

// s.listenser() listens to acitivity on s.logQueue, s.joinQueue and s.exitQueue and handles them.
func (s *server) listener() {
	for {
		select {
		// implement exit for listener
		case <-s.shutdown:
			fmt.Println("Exiting server listener...")
			return

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

		// implement exit for broadcaster
		if string(msg.body) == "\033exit" {
			break
		}

		msgPretty := formatMsg(msg)

		s.logQueue <- msg
		s.history = append(s.history, msgPretty)

		for _, cl := range s.clients {
			if cl.name == msg.from.name {
				msgPretty = []byte("\033[1m" + string(msgPretty) + "\033[0m")
			}
			_, err := cl.conn.Write(msgPretty)
			if err != nil {
				cl.conn.Close()
			}
		}
	}
	fmt.Println("Exiting server broadcastor...")
}
