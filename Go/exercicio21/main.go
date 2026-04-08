package main

import "fmt"

func verifyZero(ptr **int) *int {
	if ptr == nil {
		*ptr = new(int)
		**ptr = 22

		return *ptr
	} else {
		return *ptr
	}
}

func main() {
	var myPtr *int
	verifyZero(&myPtr)

	fmt.Println("Pointer:", &myPtr)
}
