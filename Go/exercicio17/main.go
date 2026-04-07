package main

import (
	"fmt"
	"slices"
)

func removeDuplicate(slice []string) []string {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				slice = slices.Delete(slice, j, j+1)
			}
		}
	}
	return slice
}

func main() {
	array := []string{"A", "A", "b", "c", "d"}

	fmt.Println(removeDuplicate(array))
}
