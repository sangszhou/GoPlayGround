package main

import (
	"net/rpc"
	"net"
	"net/http"
	"log"
	"learn_basic/rpc/simple/server"
	"time"
)

func main() {

	arith := new(server.Arith)

	rpc.Register(arith)

	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")

	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, nil)

	print("servering")
	time.Sleep(3600*1000)

}

