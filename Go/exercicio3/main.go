package main

import "fmt"

func Soma(num1 int, num2 int, num3 int, num4 int, num5 int) int {
	return num1 + num2 + num3 + num4 + num5
}

func Media(num1 float32, num2 float32, num3 float32, num4 float32, num5 float32) float32 {
	return (num1 + num2 + num3 + num4 + num5) / 5
}

func main() {
	fmt.Println(Soma(1, 2, 3, 4, 5))
	fmt.Println(Media(15, 20, 15, 20, 20))
}
