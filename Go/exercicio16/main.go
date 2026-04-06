package main

import "fmt"

func maxValue(slice []int) int {

	indexMaxValue := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > indexMaxValue {
			indexMaxValue = slice[i]
		}
	}
	return indexMaxValue
}

func main() {
	Varvalue := []int{1, 2, 3, 4, 5}

	fmt.Println("Max Value:", maxValue(Varvalue))
}
