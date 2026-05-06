package main

import "fmt"

type Address struct {
	Street  string
	City    string
	Country string
}

type Company struct {
	Name     string
	Location Address
}

func main() {

	MoonLab := Company{
		Name: "Moon Lab",
		Location: Address{
			Street:  "Tv. MFDOOM 15",
			City:    "Porto",
			Country: "Portugal",
		},
	}

	ptrMoonLab := &MoonLab
	ptrMoonLab.Location.Street = "Av. Liberdade 12"
	fmt.Println(MoonLab)
}
