package main

import (
	"fmt"
)

func newMatrixPointer(v ...*[]int) []int {
	pointer := new([]int)

	for _, v1 := range v {
		*pointer = append(*pointer, *v1...)
	}
	return *pointer
}

func main() {

	var values1 = []int{1, 2, 3, 4, 5}
	var values2 = []int{6, 7, 8, 9, 10}

	res := newMatrixPointer(&values1, &values2)
	fmt.Println(res)
}
