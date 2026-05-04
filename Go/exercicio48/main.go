package main

import (
	"fmt"
)

type Retangle struct {
	Height float64
	Width  float64
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

func main() {
	r := NewRetangle(-4, 2.4)
	fmt.Println(r)
}
