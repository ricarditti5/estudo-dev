package main

import "fmt"

func main() {
	var (
		twoDArray = [][]int{
			{1, 2, 89},
			{10, 20, 30, 40, 89},
			{89, 0},
		}
		targetValue int = 89
	)

	for i := range twoDArray {
		for j := range twoDArray[i] {
			if twoDArray[i][j] == targetValue {
				fmt.Println("Target value found in ", i, "position and ", j)
				break
			}
		}
	}
}
