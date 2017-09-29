package main

import (
	"net/rpc"
	"fmt"
	"log"
	"learn_basic/rpc/simple/server"
)

func main() {
	serverAddress := "localhost"

	client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &server.Args{7,8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	// Asynchronous call
	//quotient := new(Quotient)
	//divCall := client.Go("Arith.Divide", args, quotient, nil)
	//<- divCall.Done	// will be equal to divCall
	// check errors, print, etc.

}
