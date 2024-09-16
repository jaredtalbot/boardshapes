package main

import "fmt"

func main() {
	array := [5]int{1, 2, 3, 4, 5} // initialize array
	slice := array[1:4]            // create a slice of the array
	// this slice will contain elements at indices >= 1 and < 4
	fmt.Printf("array = %v\n", array)
	fmt.Printf("slice = array[1:4] = %v\n", slice)
	// Expected Output:
	// 	array = [1 2 3 4 5]
	// slice = array[1:4] = [2 3 4]
}
