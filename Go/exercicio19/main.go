package main

import "fmt"

func main() {
	var (
		value1 int8
		value2 uint32
		value3 float32
		value4 string
		value5 bool
	)

	fmt.Printf("Value 1: %v\n Value 2: %v\n Value 3: %v\n Value 4: %v\n Value 5: %v\n",
		value1, value2, value3, value4, value5)
}
