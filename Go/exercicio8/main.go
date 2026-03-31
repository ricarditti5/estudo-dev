package main

import "fmt"

func main() {

	var (
		productName  string  = "Laptop"
		productPrice float64 = 1200.50
	)
	discount := true
	quantity := 10

	fmt.Println("Product Name: "+productName, "\nPrice: ", productPrice, "\nDiscount: ", discount, "\nStock: ", quantity)
}
