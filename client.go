package main

import (
	"bufio"
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
	in   chan []byte
	out  chan []byte
	exit chan struct{}
}

func (cl *client) setup() {
	cl.setName()
	go cl.sender()
	go cl.receiver()
}

func (cl *client) setName() {
	if _, err := cl.conn.Write([]byte(welcomeMsg)); err == nil {
		if scanner := bufio.NewScanner(cl.conn); scanner.Scan() {
			if scanner.Err() != nil {
				return
			}
			cl.name = strings.TrimSpace(scanner.Text())
			cl.out <- []byte(cl.name + " has joined.")
		}
	}
	cl.exit <- struct{}{}
}

func (cl *client) sender() {
	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		cl.out <- scanner.Bytes()
	}
	cl.exit <- struct{}{}
}

func (cl *client) receiver() {
	for msg := range cl.in {
		_, err := cl.conn.Write(msg)
		if err != nil {
			break
		}
	}
	cl.exit <- struct{}{}
}
