package main

import "fmt"

func main() {
	var i int

	for i = 0; i <= 20; i++ {
		if i%3 == 0 {
			continue
		}
		fmt.Println(i)
	}
}
