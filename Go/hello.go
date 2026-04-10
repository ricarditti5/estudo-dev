package main

import "fmt"

func main() {
	name := "Ricardo"
	fmt.Println("Hello, ", name+"!")
	var i int32 = 1000
	var j int64 = int64(i) // Convert int32 to int64
	fmt.Println(j)

	var k int8 = int8(i) // Convert int32 to int8 (potential data loss)
	fmt.Println(k)       // Output: -24 (due to overflow)
}
