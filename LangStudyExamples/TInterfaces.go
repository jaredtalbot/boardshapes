package main

import (
	"fmt"
	"math"
)

type ShapeWithArea interface { //basic interface
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func main() {
	var shape ShapeWithArea = Circle{Radius: 5}
	fmt.Println(shape.Area())
	shape = Rectangle{Width: 5, Height: 4}
	fmt.Println(shape.Area())
	// Expected Output:
	// 78.53981633974483
	// 20
}
