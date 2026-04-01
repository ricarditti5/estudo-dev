package main

import "fmt"

func main() {
	idd, cidadania := 17, "Portugusa"

	if idd >= 18 && cidadania == "Portuguesa" || cidadania == "portuguesa" {
		fmt.Println("Esta pessoa pode Votar")
	} else {
		fmt.Println("Esta pessoa não pode votar")
	}
}
