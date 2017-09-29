package main

import (
	"learn_basic/rpc/recursiveCall/servers"
	"net/rpc"
	"net"
	"net/http"
	"log"
)

func main() {
	helpDesk := new(servers.Information)
	rpc.Register(helpDesk)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", "localhost:4432")

	if e != nil {
		log.Printf("Failed to start help desk ")
	}

	http.Serve(l, nil)

}
