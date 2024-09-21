package main

import (
	"fmt"
	"math/cmplx"
)

type ExampleStruct struct {
	A, B int16
}

var (
	boolEx       bool              = false                             // boolean
	stringEx     string            = "Hello, World!"                   // string
	intEx        int               = 1                                 // platform-dependent signed integer (32-bit on 32-bit systems, 64-bit on 64-bit systems)
	int8Ex       int8              = 127                               // 8-bit signed integer (signed byte)
	int16Ex      int16             = 32767                             // 16-bit signed integer
	int32Ex      int32             = 2147483647                        // 32-bit signed integer
	int64Ex      int64             = 9223372036854775807               // 64-bit signed integer
	uintEx       uint              = 18446744073709551615              // platform-dependent unsigned integer (32-bit on 32-bit systems, 64-bit on 64-bit systems)
	uint8Ex      uint8             = 255                               // 8-bit unsigned integer (unsigned byte)
	uint16Ex     uint16            = 65535                             // 16-bit unsigned integer
	uint32Ex     uint32            = 4294967295                        // 32-bit unsigned integer
	uint64Ex     uint64            = 18446744073709551615              // 64-bit unsigned integer
	float32Ex    float32           = 3.4e+38                           // 32-bit floating point number
	float64Ex    float64           = +1.7e+308                         // 64-bit floating point number
	complex64Ex  complex64         = 3 + 2i                            // 64-bit complex number
	complex128Ex complex128        = cmplx.Sqrt(64 + 9i)               // 128-bit complex number
	structEx     ExampleStruct     = ExampleStruct{A: 123, B: 321}     // struct
	interfaceEx  any               = "Anything"                        // any type; can be asserted to a concrete type later. this is an alias for interface{}
	arrayEx      [3]int            = [3]int{1, 2, 3}                   // array
	sliceEx      []int             = arrayEx[0:1]                      // slice
	mapEx        map[string]string = map[string]string{"Key": "Value"} // map
	channelEx    chan string       = make(chan string)                 // channel
)

func main() {
	fmt.Printf("Type: %-20T Value: %v\n", boolEx, boolEx)
	fmt.Printf("Type: %-20T Value: %v\n", stringEx, stringEx)
	fmt.Printf("Type: %-20T Value: %v\n", intEx, intEx)
	fmt.Printf("Type: %-20T Value: %v\n", int8Ex, int8Ex)
	fmt.Printf("Type: %-20T Value: %v\n", int16Ex, int16Ex)
	fmt.Printf("Type: %-20T Value: %v\n", int32Ex, int32Ex)
	fmt.Printf("Type: %-20T Value: %v\n", int64Ex, int64Ex)
	fmt.Printf("Type: %-20T Value: %v\n", uintEx, uintEx)
	fmt.Printf("Type: %-20T Value: %v\n", uint8Ex, uint8Ex)
	fmt.Printf("Type: %-20T Value: %v\n", uint16Ex, uint16Ex)
	fmt.Printf("Type: %-20T Value: %v\n", uint32Ex, uint32Ex)
	fmt.Printf("Type: %-20T Value: %v\n", uint64Ex, uint64Ex)
	fmt.Printf("Type: %-20T Value: %v\n", float32Ex, float32Ex)
	fmt.Printf("Type: %-20T Value: %v\n", float64Ex, float64Ex)
	fmt.Printf("Type: %-20T Value: %v\n", complex64Ex, complex64Ex)
	fmt.Printf("Type: %-20T Value: %v\n", complex128Ex, complex128Ex)
	fmt.Printf("Type: %-20T Value: %v\n", structEx, structEx)
	fmt.Printf("Type: %-20T Value: %v\n", interfaceEx, interfaceEx)
	fmt.Printf("Type: %-20T Value: %v\n", arrayEx, arrayEx)
	fmt.Printf("Type: %-20T Value: %v\n", sliceEx, sliceEx)
	fmt.Printf("Type: %-20T Value: %v\n", mapEx, mapEx)
	fmt.Printf("Type: %-20T Value: %v\n", channelEx, channelEx)
}
