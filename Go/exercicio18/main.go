package main

import "fmt"

func appendSlices(slice1 []int, slice2 []int) []int {
	appendedSlice := append(slice1, slice2...)

	return appendedSlice
}

func main() {

	fatia1 := []int{1, 2, 3, 4}
	fatia2 := []int{5, 6, 7, 8}

	fmt.Println(appendSlices(fatia1, fatia2))
}
