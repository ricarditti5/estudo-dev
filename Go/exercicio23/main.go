package main

import (
	"fmt"
	"math"
)

func main() {

	const PI float64 = math.Pi
	var r float64 = 4
	resAreaCircle := PI * r * r

	fmt.Printf("The area of circle is: %.3v cm^2", resAreaCircle)
}
