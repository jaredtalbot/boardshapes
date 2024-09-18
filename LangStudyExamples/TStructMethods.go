package main

import "fmt"

type Animal struct { // type Animal
	Species string
}

func (a *Animal) SetSpeciesToDog() string { // method that extends animal
	a.Species = "dog"
	return a.Species
}

func main() {
	animal := Animal{"cat"}
	fmt.Printf("Animal Species before method call: %s\n", animal.Species)
	animal.SetSpeciesToDog() // this method will have the effect of changing animal.Species
	fmt.Printf("Animal Species after method call: %s\n", animal.Species)
	// Expected Output:
	// Animal Species before method call: cat
	// Animal Species after method call: dog
}
