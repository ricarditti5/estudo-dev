package main

import "fmt"

func multplesNumber(x int) {
	if x%3 == 0 {
		fmt.Printf("fizz")
	} else if x%5 == 0 {
		fmt.Printf("bizz")
	} else if x%3 == 0 && x%5 == 0 {
		fmt.Printf("FizzBuzz")
	} else {
		fmt.Println("Something wrong")
	}
}

func main() {
	var num int

	fmt.Println("Type a random number: ")
	fmt.Scanf("%d", &num)

	multplesNumber(num)
}
