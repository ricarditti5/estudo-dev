package main

import "fmt"

func main() {
	var (
		num1 int8   = 8
		num2 int16  = 80
		num3 int32  = 800
		num4 int64  = 8000
		num5 uint8  = 80
		num6 uint16 = 8000
		num7 uint32 = 800000
		num8 uint64 = 80000000000
	)

	fmt.Printf("%T %v - %T %v - %T %v - %T %v - %T %v - %T %v - %T %v - %T %v\n",
		num1, num1, num2, num2, num3, num3, num4, num4,
		num5, num5, num6, num6, num7, num7, num8, num8)
}
