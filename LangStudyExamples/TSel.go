package main

import (
	"fmt"
	"time"
)

var (
	a, b, c, d, e, f = 1, 2, 3, 4, 5, 6
)

func main() {
	// if: run code if expression is true
	if d > a {
		fmt.Println("d > a")
	}
	// if else: fall back to else if the expression is not true
	if (e < b) && (c >= b) {
		fmt.Println("(e < b) && (c >= b)")
	} else {
		fmt.Println("(e >= b) || (c < b)")
	}
	// nested if else
	if a != b {
		fmt.Println("(a != b)")
	} else {
		if (d == e) || (e != f) {
			fmt.Println("(a == b)&& ((d == e) || (e != f))")
		}
	}

	// switch case: run code for whichever case that matches, otherwise default
	switch d {
	case 1:
		fmt.Println("d is 1")
	case 2, 3, 4:
		fmt.Println("d is 2, 3, or 4")
	default:
		fmt.Println("d is some other number")
	}

	chan1, chan2 := make(chan string), make(chan string)

	go func() {
		time.Sleep(5 * time.Second)
		chan1 <- "foo"
	}()

	go func() {
		time.Sleep(3 * time.Second)
		chan2 <- "bar"
	}()

	// select: waits until data can be sent/received through a case's channel, and then runs the code for that case
	select {
	case msg := <-chan1:
		fmt.Printf("chan1 said: %s\n", msg)
	case msg := <-chan2:
		fmt.Printf("chan2 said: %s\n", msg)
	}
	// (only chan2 will be received here)
}

// Expected Output:
// d > a
// (e >= b) || (c < b)
// (a != b)
// d is 2, 3, or 4
// chan2 said: bar
