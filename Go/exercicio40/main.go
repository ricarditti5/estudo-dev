package main

import "fmt"

func pointerShow(v1 *int, v2 *int) {
	*v1 = 30
	*v2 = 50
}

func main() {

	var1 := 10
	var2 := 40
	fmt.Println(var1, var2)

	pointerShow(&var2, &var1)
	fmt.Println(var2, var1)
}
