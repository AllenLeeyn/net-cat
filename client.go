package main

import (
	"bufio"
	"net"
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
	from chan []byte
	exit chan struct{}
}

func (cl *client) getFrom() {
	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		cl.from <- []byte(scanner.Text() + "\n")
	}
	cl.exit <- struct{}{}
}
