package main

import (
	"fmt"
	"math"
)

type Circle struct {
	diameter float64
}

func (t Circle) Circumference() float64 {
	return t.diameter * math.Pi
}

func main() {

	a := Circle{diameter: 2.4}

	fmt.Printf("The circunference is: %.2f\n", a.Circumference())

}
