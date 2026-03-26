package main

import "fmt"

var variavelExt string = "Variável externa"

func main() {
	var variavelInt string = "Variável interna"

	fmt.Println(variavelExt+"\n", variavelInt)
}
