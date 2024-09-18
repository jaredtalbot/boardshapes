package main

import (
	"fmt"
)

func main() {
	var m map[string]int     // map with keytype string and valuetype int
	m = make(map[string]int) // initializing a map
	m["route"] = 66          // we can set the key route to value 66
	m["area"] = 51           // also set the key area to value 51
	r := m["route"]          // retrieve the value at route
	fmt.Printf("route = %d\n", r)
	a, exists := m["area"] // retrieve the value at area; the second value will be true since there is a value at area
	fmt.Printf("area = %d (exists = %v)\n", a, exists)
	o, exists := m["order"] // retrieve the value at order; the second value will be false since there is no value at order
	fmt.Printf("order = %d (exists = %v)\n", o, exists)
	// Expected Output:
	// route = 66
	// area = 51 (exists = true)
	// order = 0 (exists = false)
}
