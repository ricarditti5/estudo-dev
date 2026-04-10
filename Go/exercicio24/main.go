package main

import (
	"fmt"
)

func main() {
	const (
		StatusOk = iota
		NotFound
		InternalServerError
	)

	fmt.Println(StatusOk, " ", NotFound, " ", InternalServerError)
}
