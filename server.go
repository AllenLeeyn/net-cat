package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const maxConn = 2

type server struct {
	clients   []*client
	msgQueue  chan message
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

	s.msgQueue <- message{from: "server", body: []byte("Listening on port " + portNum + "\n")}

	go s.listener()
	go s.broadcaster()

	for {
		conn, err := server.Accept()
		if err != nil || len(s.clients) >= maxConn {
			conn.Write([]byte("Server full. Try again later.\n"))
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
		exit: make(chan struct{})}

	s.addClient(cl)

	<-cl.exit
	s.removeClient(cl)
}

func (s *server) addClient(cl *client) {
	_, err := cl.conn.Write([]byte(welcomeMsg))
	if err == nil {
		scanner := bufio.NewScanner(cl.conn)
		for scanner.Scan() {
			cl.name = strings.TrimSpace(scanner.Text())
			if cl.name == "" {
				cl.conn.Write([]byte("Nothing entered. Try again: "))
				continue
			}
			for _, others := range s.clients {
				if others.name == cl.name {
					cl.conn.Write([]byte("Name taken. Try again: "))
					continue
				}
			}
			s.clients = append(s.clients, cl)
			s.msgQueue <- message{from: "server",
				body: []byte(cl.name + " has joined the chat.\n")}
			go cl.getFrom()
			return
		}
	}
	cl.exit <- struct{}{}
	s.removeClient(cl)
}

func (s *server) removeClient(cl *client) {

	for i, c := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	cl.conn.Close()
	close(cl.from)
	close(cl.exit)
	s.msgQueue <- message{from: "server",
		body: []byte(cl.name + " has leaved the chat.\n")}
}

func (s *server) listener() {
	fmt.Println("listening")
	for {
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

func (s *server) broadcaster() {
	fmt.Println("broadcasting")
	for {
		msg := <-s.msgQueue
		col := cols["red"]
		if msg.from != "server" {
			col = cols["blue"]
		}
		msgPretty := formatMsg(msg.from, string(msg.body), col)
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
