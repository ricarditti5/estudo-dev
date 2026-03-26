package main

import (
	stringutils "exercicio2/package"
	"fmt"
)

func main() {
	var input string = "Ola"

	fmt.Println("String Original " + input)
	fmt.Println("String Invertida " + stringutils.ReverseString(input))
}
