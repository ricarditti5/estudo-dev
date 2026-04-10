package main

import "fmt"

func main() {
	const (
		c = 22.3
		f = 44.5
	)

	resConvertionF := (c * 1.8) + 32
	resConvertionC := (f - 32) * 5 / 9

	fmt.Println("Celsius to Fahrenheit:", resConvertionF)
	fmt.Printf("Fahrenheit to Celsius: %.3v", resConvertionC)
}
