package main

import (
	"fmt"
)

type Retangle struct {
	Height float64
	Width  float64
}
type Area struct {
	Retangle
}

func NewRetangle(height, width float64) *Retangle {
	if height < 0 || width < 0 {
		return nil
	}
	return &Retangle{
		Height: height,
		Width:  width,
	}
}

func (t Retangle) CalculateArea() float64 {
	return t.Height * t.Width
}
func (t Retangle) CalculatePerimeter() float64 {
	return 2 + (t.Height + t.Width)
}

func main() {
	r := NewRetangle(4, 2.4)
	a := Area{
		Retangle{
			Height: 3,
			Width:  2,
		},
	}

	fmt.Println(r)
	fmt.Println(a)

	fmt.Println("The area is: ", a.CalculatePerimeter())
	fmt.Println("The Perimeter is: ", a.CalculateArea())

}
