package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan string) // make a channel for strings
	go func() {                  // run anonymous goroutine
		message := <-channel // receive string from channel
		fmt.Printf("received: %s", message)
	}()
	fmt.Println("sending to goroutine...")
	channel <- "hello there!" // send string through channel
	// Expected Output:
	// sending to goroutine...
	// received: hello there!
	time.Sleep(2 * time.Second)
}
