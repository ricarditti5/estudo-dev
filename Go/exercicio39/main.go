package main

import "fmt"

func calculator(initialValue float64) (func(float64) float64, func(float64) float64, func(float64) float64, func(float64) float64) {
	res := initialValue

	add := func(n float64) float64 {
		res += n
		return res
	}

	subtract := func(n float64) float64 {
		res -= n
		return res
	}

	multiply := func(n float64) float64 {
		res *= n
		return res
	}

	divide := func(n float64) float64 {
		res /= n
		return res
	}
	return add, subtract, multiply, divide
}

func main() {
	add, subtract, multiply, divide := calculator(10)

	fmt.Println(add(5))      // Output: 15
	fmt.Println(subtract(3)) // Output: 12
	fmt.Println(multiply(2)) // Output: 24
	fmt.Println(divide(4))   // Output: 6
}
