package main

import (
	"fmt"
)

func verifyPointer(p *int) int {
	if p == nil {
		fmt.Println("empty pointer")
	}

	return *p
}

func main() {
	v := 10

	fmt.Println(verifyPointer(&v))
}
