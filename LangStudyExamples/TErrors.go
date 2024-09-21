package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) { // this method has an error as a second return value
	if b == 0 { // we can't divide by zero
		return 0, errors.New("divide: division by zero not allowed") // return default value and error
	}
	return a / b, nil // return quotient and no error
}

func main() {
	a := 4
	b := 0

	d, err := divide(a, b)
	if err != nil { // error check
		fmt.Println(err) // there is an error, print it...
		return           // ...and return
	}

	// at this point we have asserted there was no error
	fmt.Printf("Result: %d\n", d)
}

//Expected Output:
//"divide: division by zero not allowed"
