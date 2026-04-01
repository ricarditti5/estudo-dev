package main

import (
	"fmt"
	"math"
)

func main() {
	var (
		nota1 float64 = 16.34
		nota2 float64 = 18.621
	)
	res := nota1 + nota2/2

	fmt.Printf("A Media das notas é: %.1f", res)

	x := math.Abs(res)
	fmt.Printf("\nA Media Absoluta das notas é: %.1f", x)

}
