package main

import "fmt"

func numberCompare(x int, y int) {
	if x > y {
		fmt.Printf("Number %d is bigger than %d", x, y)
	} else if x < y {
		fmt.Printf("Number %d is smaller than %d", x, y)
	} else if x == y {
		fmt.Printf("Number %d is iqual %d", x, y)
	} else {
		fmt.Printf("Something wrong")
	}
}

func main() {

	var (
		num1 int
		num2 int
	)

	fmt.Println("\nType int number:")
	fmt.Scanf("%d", &num1)
	fmt.Scanf("%d", &num2)

	numberCompare(num1, num2)
}
