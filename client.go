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
	"     `-'       `--'\n" +
	"[ENTER YOUR NAME]:"

type client struct {
	conn net.Conn
	name string
}

func (s *server) addClient(cl *client) {
	_, err := cl.conn.Write([]byte(welcomeMsg))
	if err == nil {
		scanner := bufio.NewScanner(cl.conn)
		for scanner.Scan() {
			cl.name = strings.TrimSpace(scanner.Text())
			if cl.name == "" || !isValidEntry(cl.name) || s.isNameTaken(cl.name) {
				cl.conn.Write([]byte("Invalid entry/ name taken. Try again: "))
				continue
			}
			if err := s.printHistory(cl); err != nil {
				break
			}
			s.logQueue <- message{from: "server",
				body: []byte(cl.conn.RemoteAddr().String() + " set name to " + cl.name)}
			s.joinQueue <- cl
			return
		}
	}
	s.logQueue <- message{from: "server",
		body: []byte("Unable to connect " + cl.conn.RemoteAddr().String())}
	cl.conn.Close()
}

func isValidEntry(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	for _, rn := range entry {
		if !(rn >= 32 || rn <= 126 ||
			rn == 'å' || rn == 'ä' || rn == 'ö' ||
			rn == 'Å' || rn == 'Ä' || rn == 'Ö') {
			return false
		}
	}
	return true
}

func (s *server) isNameTaken(name string) bool {
	for _, cl := range s.clients {
		if cl.name == name {
			return true
		}
	}
	return false
}

func (s *server) printHistory(cl *client) error {
	for _, msg := range s.history {
		_, err := cl.conn.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to write history")
		}
	}
	return nil
}

func (s *server) removeClient(cl *client) {
	for i, c := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			s.msgQueue <- message{from: "server",
				body: []byte(cl.name + " has left the chat.")}
		}
	}
	s.logQueue <- message{from: "server",
		body: []byte("Close connection with " + cl.conn.RemoteAddr().String())}
	cl.conn.Close()
}
