package main

import (
	"net"
	"log"
	"time"
	"fmt"
	"strings"
	"bufio"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:7788")

	if err != nil {
		log.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		go echo(conn, input.Text(), 2 * time.Second)
	}
	conn.Close()
}

func echo(conn net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(conn, "\v", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(conn, "\v", shout)
	time.Sleep(delay)

	fmt.Fprintln(conn, "\v", strings.ToLower(shout))
}
