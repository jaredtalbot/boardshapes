package main

import "fmt"

func add(x, y int) int {
	return x + y
}

func addprint(x, y int) {
	fmt.Printf("%d + %d = %d\n", x, y, x+y)
}

func div(dividend, divisor int) (quotient, remainder int) {
	return dividend / divisor, dividend % divisor
}

func swap(x, y *int) {
	*x, *y = *y, *x
}

func main() {
	// function with no return value
	addprint(23, 323)

	// function with single return value
	a, b := 42, 23
	i := add(a, b)
	fmt.Printf("%d + %d = %d\n", a, b, i)

	// function with multiple return values
	dd, dr := div(a, b)
	fmt.Printf("%d / %d = %d with remainder %d\n", a, b, dd, dr)

	x := 100
	y := 200

	fmt.Printf("x = %d, y = %d\n", x, y)
	// function that accepts pointers and modifies the referenced values
	swap(&x, &y)
	fmt.Printf("x = %d, y = %d\n", x, y)

	// Expected Output:
	// 23 + 323 = 346
	// 42 + 23 = 65
	// 42 / 23 = 1 with remainder 19
	// x = 100, y = 200
	// x = 200, y = 100
}
