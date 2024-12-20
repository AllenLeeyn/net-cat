package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
	"     `-'       `--'\n"

type client struct {
	conn  net.Conn
	name  string
	color string
}

// s.processClient() sends a welcome message and checks if name is valid
// before adding connecting client to s.joinQueue.
func (s *server) processClient(cl *client) {
	_, err := cl.conn.Write([]byte(welcomeMsg))
	if err == nil {
		if s.setClientName(cl) {
			return
		}
	}
	s.logQueue <- message{from: s.self,
		body: []byte("Unable to connect " + cl.conn.RemoteAddr().String())}
	cl.conn.Close()
}

func (s *server) setClientName(cl *client) bool {
	_, err := cl.conn.Write([]byte("[ENTER YOUR NAME]:"))
	if err == nil {
		oldName := cl.name
		cl.name = ""

		scanner := bufio.NewScanner(cl.conn)
		for scanner.Scan() {
			newName := strings.TrimSpace(scanner.Text())
			if newName == "" || !isValidEntry(newName) || s.isNameTaken(newName) {
				cl.conn.Write([]byte("Invalid entry/ name taken. Try again: "))
				continue
			}
			cl.name = newName
			s.logQueue <- message{from: s.self,
				body: []byte(cl.conn.RemoteAddr().String() + " set name to " + cl.name)}

			if oldName == "" {
				s.joinQueue <- cl
			} else {
				s.msgQueue <- message{from: cl,
					body: []byte(oldName + " changed name to " + cl.name)}
			}
			return true
		}
	}
	return false
}

// s.isNameTaken() checks if requested name is used by any registered clients.
func (s *server) isNameTaken(name string) bool {
	for _, cl := range s.clients {
		if cl.name == name {
			return true
		}
	}
	return false
}

// s.printHistory() prints history to the new client.
func (s *server) printHistory(cl *client) error {
	for _, msg := range s.history {
		_, err := cl.conn.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to write history")
		}
	}
	return nil
}

// s.removeClient() removes client from the client list,
// add exit message to s,msgQueue, logs activity and close conn.
func (s *server) removeClient(cl *client) {
	for i, c := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			s.msgQueue <- message{from: s.self,
				body: []byte(cl.name + " has left the chat.")}
		}
	}
	s.logQueue <- message{from: s.self,
		body: []byte("Close connection with " + cl.conn.RemoteAddr().String())}
	cl.conn.Close()
}
