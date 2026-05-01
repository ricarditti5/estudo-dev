package main

import (
	"fmt"
)

func addPrefix(prefix ...[]string) []string {

	var newPrefix []string

	for _, s := range prefix {
		newPrefix = append(newPrefix, s...)
	}

	return newPrefix

}

func main() {
	a := []string{"aaa"}
	b := []string{"bbbb", "ccccccc"}

	fmt.Println(addPrefix(a, b))
}
