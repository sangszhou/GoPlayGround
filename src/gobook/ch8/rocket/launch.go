package main

import (
	"time"
	"log"
)

func main() {
	abort := make(chan interface{})


	select {
	case <- time.After(10 * time.Second):
		// do nothing
	case <- abort:
		log.Println("launch aborted")
		return
	}

	log.Println("launch")

}
