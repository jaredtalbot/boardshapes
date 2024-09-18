package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan string)
	go func() { // run anonymous goroutine
		message := <-channel
		fmt.Printf("received: %s", message)
	}()
	fmt.Println("sending to goroutine...")
	channel <- "hello there!"
	// Expected Output:
	// sending to goroutine...
	// received: hello there!
	time.Sleep(2 * time.Second)
}
