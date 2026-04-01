package main

import "fmt"

func main() {

	var inputString string

	fmt.Println("Diga uma string para executar operações:")
	fmt.Scan(&inputString)

	//Show the length of a string
	fmt.Println("\nThe Length of the string is: ", len(inputString))
	//Access and print the first 3 letters of the string
	fmt.Println("\nThe First 3 letters of the string is: ", inputString[0:3])

	secondString := ""
	fmt.Println("\nDiga uma Outra string para executar operações:")
	fmt.Scan(&secondString)

	//Concat the inputString and secondString
	fmt.Println("\nThose strings are concat: " + inputString + "-" + secondString)
}
