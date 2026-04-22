package main

import (
	"fmt"
	"math/rand"
)

func randomNumbers() int {
	return rand.Intn(4)
}

func main() {
	num := randomNumbers()

	switch num {
	case 1:
		fmt.Println("Green")
	case 2:
		fmt.Println("Yellow")
	case 3:
		fmt.Println("Red")
	default:
		fmt.Println("Sinal Inválido")
	}
}
