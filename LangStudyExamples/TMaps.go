package main

import "fmt"

func main() {
	var m map[string]int     // map with keytype string and valuetype int
	m = make(map[string]int) // initializing a map
	m["route"] = 66          // we can set the key route to value 66
	r := m["route"]          //retrieve the value at route
	fmt.Println(r)
	// Expected Output:
	// 66
}
