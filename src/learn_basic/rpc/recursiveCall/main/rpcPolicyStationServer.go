package main

import (
	"learn_basic/rpc/recursiveCall/servers"
	"net/rpc"
	"net"
	"log"
	"net/http"
)

func main() {
	policeMan := new(servers.PoliceMan)
	rpc.Register(policeMan)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":9110")

	if e != nil {
		log.Printf("Failed to start police station ")
	}

	http.Serve(l, nil)
}
