package main

import "fmt"

func main() {
	nota := 2

	switch {
	case nota >= 0 && nota <= 56:
		fmt.Println("F")
	case nota >= 60 && nota <= 69:
		fmt.Println("D")
	case nota >= 70 && nota <= 79:
		fmt.Println("C")
	case nota >= 80 && nota <= 89:
		fmt.Println("B")
	case nota >= 90 && nota <= 100:
		fmt.Println("A")
	}
}
