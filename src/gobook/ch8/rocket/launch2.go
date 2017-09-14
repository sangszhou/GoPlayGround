package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print("hello")

	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C    // receive from the ticker's channel
	ticker.Stop() // cause the ticker's goroutine to terminate

}
