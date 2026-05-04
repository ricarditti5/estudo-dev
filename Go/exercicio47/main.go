package main

import "fmt"

type Student struct {
	Name    string
	ID      int
	Grades  []float64
	Address Address
}

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

func main() {
	p := Student{
		Name:   "Ricardo",
		ID:     1,
		Grades: []float64{20, 19, 18, 16, 17},
		Address: Address{
			Street:  "Travessa de São Bento 18",
			City:    "Barcelos",
			State:   "Braga",
			ZipCode: "4750-268",
		},
	}
	fmt.Println("The address is: ", p.Address)
}
