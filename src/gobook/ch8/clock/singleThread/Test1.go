package main

import (
	"net"
	"log"
	"io"
	"time"
)

//使用 nc -c localhost 7788 来连接
//只有一个 client 能够收到时间，其他的 client 都不受到
//client 退出是，报错是 write tcp 127.0.0.1:7788->127.0.0.1:60935: write: broken pipe

func main() {
	listener, err := net.Listen("tcp", "localhost:7788")
	if err != nil {
		log.Println("%v", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
		}

		handleConn(conn)
	}

}

func handleConn(conn net.Conn)  {
	for {
		_, err := io.WriteString(conn, time.Now().Format("11:12:03\n"))
		if err != nil {
			log.Println("%v", err)
			return
		}

		time.Sleep(1 * time.Second)
	}
}
