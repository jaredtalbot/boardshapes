package main

import (
	"fmt"
)

func main() {
	a := 4
	b := 0
	d := a / b // division by zero
	fmt.Printf("Result: %d\n", d)
	// Expected Error Output:
	// panic: runtime error: integer divide by zero
	// (stack trace)
}
