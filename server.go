package main

import (
	"net"
)

const maxConn = 10

type server struct {
	clients   []*client
	msgQueue  chan []byte
	exitQueue chan *client
	shutdown  chan struct{}
}

func (s *server) start(portNum string) {
	listener, err := net.Listen("tcp4", portNum)
	check(err)

	logger("server", "Listening: "+portNum, cols["blue"])

	go s.broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil || len(s.clients) >= maxConn {
			continue
		}
		go s.handlerConnection(conn)
	}

	//<-s.shutdown
}

func (s *server) handlerConnection(conn net.Conn) {
	cl := &client{conn: conn,
		in:  make(chan []byte),
		out: make(chan []byte)}

	cl.setup()
	s.clients = append(s.clients, cl)

}

func (s *server) broadcaster() {
	for {
		select {
		case msg := <-s.msgQueue:
			for _, cl := range s.clients {
				if cl != nil {
					select {
					case cl.in <- msg:
					default:
					}
				}
			}
		}
	}
}
