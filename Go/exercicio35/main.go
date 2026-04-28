package main

import (
	"fmt"
	"math"
)

func main() {

	var (
		num  int
		rage int
	)

	fmt.Println("Type a number: ")
	fmt.Scanln(&num)

	rage = int(math.Sqrt(float64(num)))
	for i := 2; i <= rage; i++ {
		if num%i == 0 {
			fmt.Println(num, "Isn't a prime number")
			break
		}
	}
}
