package main

import "fmt"

type MyInt int32 // the underlying type is just an integer

func (i MyInt) Square() MyInt { // declaring a method on this integer type
	return i * i
}

func main() {
	i := MyInt(5)
	fmt.Println(i.Square()) // Expected Output: 25
}

// Expected Output:
// 25
