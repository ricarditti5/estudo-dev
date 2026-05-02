package main

import "fmt"

type student struct {
	name      string
	age       int
	isStudent bool
}

func main() {
	s := student{name: "Ricardo", age: 17, isStudent: true}
	p := &s //that point to address of s instance
	p.age = 18

	fmt.Println(s)

	p.isStudent = false
	p.age = 17

	fmt.Println(s)
}
