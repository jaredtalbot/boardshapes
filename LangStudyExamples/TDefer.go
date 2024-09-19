package main

import "fmt"

func main() {
	defer fmt.Println("This will happen fourth!")
	defer fmt.Println("This will happen third!")
	fmt.Println("This will happen first!")
	fmt.Println("This will happen second!")
	// Expected Output:
	// This will happen first!
	// This will happen second!
	// This will happen third!
	// This will happen fourth!
}
