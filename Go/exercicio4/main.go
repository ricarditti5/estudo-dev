package main

import "fmt"

var variavel string = "Variável externa"

func main() {
	var variavel string = "Variável interna"

	fmt.Println(variavel+"\n", variavel)
}
