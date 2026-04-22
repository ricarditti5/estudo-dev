package main

import (
	"fmt"
	"math"
	"strconv"
)

func MenuCalculator() {
	fmt.Println("Select an operation")
	fmt.Println("\n1-addiction + ")
	fmt.Println("\n2-subtraction - ")
	fmt.Println("\n3-multiplication * ")
	fmt.Println("\n4-division / ")
}

func addiction(n1 float64, n2 float64) int {
	res := int(n1 + n2)
	return res
}

func subtraction(n1 float64, n2 float64) int {
	res := int(n1 - n2)
	return res
}

func multiplication(n1 float64, n2 float64) int {
	res := int(n1 * n2)
	return res
}

func division(n1 float64, n2 float64) float64 {
	res := n1 / n2
	roundedRes := math.Round(float64(res)*100) / 100
	return roundedRes
}
func main() {
	var (
		op   float32
		val1 float64
		val2 float64
	)
	MenuCalculator()
	fmt.Scanln(&op)

	fmt.Println("\nType two numbers to operate:")
	fmt.Scanln(&val1, &val2)

	switch op {
	case 1:
		fmt.Printf(strconv.Itoa(addiction(val1, val2)))
	case 2:
		fmt.Printf(strconv.Itoa(subtraction(val1, val2)))
	case 3:
		fmt.Printf(strconv.Itoa(multiplication(val1, val2)))
	case 4:
		fmt.Printf("%v", division(val1, val2))
	default:
		fmt.Print("Invalid operation\n" +
			"Type the right numbers and choose the correct option.\n")
	}
}
