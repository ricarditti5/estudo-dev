package main

import "fmt"

func main() {
	numbers := []int{1, 1, 3, 7, 5, 10, 7, 2, 8}

	for _, num := range numbers {
		if num%2 == 0 {
			break
		}
		fmt.Println(num)
	}
}
