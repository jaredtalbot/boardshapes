package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide: division by zero not allowed")
	}
	return a / b, nil
}

func main() {
	a := 4
	b := 0
	d, err := divide(a, b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %d\n", d)
	}
}
