package main

import (
	"fmt"
	"time"
)

func main() {
	a, b := 1, 0

	// for loop: same as other languages
	for i := 0; i < 5; i++ {
		a = a * 2
		fmt.Printf("Value of a is %d\n", a)
	}
	fmt.Println()

	// while loop: it's a for loop without the initializer or post-loop statement
	for b < 5 {
		b = b + 1
		fmt.Printf("Value of b is %d\n", b)
	}
	fmt.Println()

	myArray := [3]int{62, 74, 23}

	// for range over array: iterates through every element of the array, giving you the index and value
	for i, v := range myArray {
		fmt.Printf("myArray[%d] = %d\n", i, v)
	}
	fmt.Println()

	myMap := map[string]float64{
		"PI": 3.14159,
		"K":  8.98755e9,
		"G":  6.67430e-11,
	}

	// for range over map: iterates through every element of the map, giving you the key and value
	for k, v := range myMap {
		fmt.Printf("myMap[%s] = %E\n", k, v)
	}
	fmt.Println()

	c := make(chan string, 3)
	go func() {
		for i := 3; i > 0; i-- {
			c <- fmt.Sprint(i)
			time.Sleep(1 * time.Second)
		}
		c <- "Liftoff"
		time.Sleep(3 * time.Second)
		close(c) // close channel
	}()

	// for range over channel: keep receiving data from the channel until it is closed
	for v := range c {
		fmt.Printf("received from c: %s\n", v)
	}
	fmt.Println("c has been closed")
}
