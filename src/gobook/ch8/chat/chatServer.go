package main

import (
	"net"
	"log"
	"bufio"
	"fmt"
)

// an outgoing channel
type client chan<- string

var (
	message = make(chan string)
	// 注意这是嵌套的 chan
	enter = make(chan client)
	leave = make(chan client)
)


func main() {

	listener, err := net.Listen("tcp", "localhost:9900")
	if err != nil {
		log.Println(err)
		return
	}

	go broadcast()

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConn(client)
	}
}

func handleConn(client net.Conn) {
	ch := make(chan string)

	// send outgoing message to client
	go func(client net.Conn, ch2 chan string) {
		for msg := range ch2 {
			fmt.Fprintln(client, msg)
		}
	}(client, ch)

	who := client.RemoteAddr().String()
	ch <- "You are " + who

	enter <- ch

	message <- who + " has arrived"
	// 这一点需要注意，enter 是 chan chan<-string 类型的, ch 是 chan string 类型的
	// 所以 ch 即便没有保护任何 string, 也是可以赋值给 enter 的

	input := bufio.NewScanner(client)
	for input.Scan() {
		message <- who + ": " + input.Text()
	}


	leave <- ch
	message <- who + " has left"
	client.Close()

}

func broadcast()  {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-message:
			// 推送消息到所有的 client
			for cli := range clients {
				cli <- msg
			}

		case cli := <-enter:
			clients[cli] = true

		case cli := <-leave:
			delete(clients, cli)
			close(cli)
		}
	}
	
}
