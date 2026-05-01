package main

import (
	"fmt"
)

func avarage(num ...float64) float64 {
	var total float64

	if num == nil || len(num) == 0 {
		return 0.0
	}
	for _, nums := range num {
		total += nums
	}
	return total / float64(len(num))
}

func main() {

	a, b, c := 15.4, 14.7, 18.5

	avarage()
	fmt.Println("The avarage classification is:")
	fmt.Println(avarage(a, b, c))
}
