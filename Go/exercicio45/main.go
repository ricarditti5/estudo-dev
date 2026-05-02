package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	debug.SetGCPercent(50)

	data := make([]int, 8000)
	fmt.Println(data)

	data[0] = 40
	fmt.Println(data)

}
