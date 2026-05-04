package main

import "fmt"

type Student struct {
	Name   string
	ID     int
	Grades []float64
}

func main() {
	s := Student{
		Name:   "Ricardo",
		ID:     1,
		Grades: []float64{18, 20, 17, 14, 15},
	}
	fmt.Println(s)
}
