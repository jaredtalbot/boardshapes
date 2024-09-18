package main

import (
	"fmt"
	"math"
)

type ShapeWithArea interface { //basic interface
	Area() float64
}

type Rectangle struct { // note: we do not explicitly say here that this implements ShapeWithArea
	Width, Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 { // implement the Area method of the ShapeWithArea interface
	return r.Width * r.Height
}

func (c Circle) Area() float64 { // implement the Area method of the ShapeWithArea interface
	return math.Pi * c.Radius * c.Radius
}

func main() {
	var shape ShapeWithArea = Circle{Radius: 5} // create a ShapeWithArea variable, assign a Circle to it
	fmt.Println(shape.Area())
	shape = Rectangle{Width: 5, Height: 4} // we can reuse this ShapeWithArea variable but with a Rectangle instead
	fmt.Println(shape.Area())
	// Expected Output:
	// 78.53981633974483
	// 20
}
