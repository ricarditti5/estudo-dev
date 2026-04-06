package main

import "fmt"

func reverseSlice(slice []int) []int {

	revertedSlice := make([]int, 0, len(slice))

	for i := len(slice) - 1; i >= 0; i-- {
		revertedSlice = append(revertedSlice, slice[i])
	}
	return revertedSlice
}

func main() {
	sliceValue := []int{1, 2, 3, 4, 5}

	fmt.Println("Inverted Values:", reverseSlice(sliceValue))
}
