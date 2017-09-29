package main

import (
	"net/rpc"
	"log"
	"learn_basic/rpc/recursiveCall/servers"
)

func main() {
	serverAddress := "localhost"

	client, err := rpc.DialHTTP("tcp", serverAddress + ":4432")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	//request := &servers.RequestArgs{
	//	Content: "direction of sjtu",
	//}

	//request := &servers.RequestArgs{
	//	Content: "911 are you there",
	//}

	request := &servers.RequestArgs{
		Content: "110 喂，在吗",
	}

	/**
		***
		这个地方不能是指针，必须是实体
	 */
	var reply servers.ReplyArgs

	err = client.Call("Information.Call", request, &reply)

	if err != nil {
		log.Print("Failed to call information", err)
	} else {
		log.Print("reply: " + reply.Content)
	}




}
