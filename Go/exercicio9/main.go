package main

import "fmt"

var variavelExt string = "SOU UMA VARIAVEL EXETERNA"

func main() {
	variavelInt := "sou uma variavel interna"

	fmt.Println("Eu " + variavelExt)
	fmt.Println("Eu " + variavelInt)
}
