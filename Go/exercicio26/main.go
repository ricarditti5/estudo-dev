package main

import "fmt"

func leapYear(y int) {
	if y%4 == 0 && y%100 != 0 || y%400 == 0 {
		fmt.Printf("%d is a Leap Year", y)
	} else {
		fmt.Printf("%d isn't a Leap Year", y)
	}
}

func main() {
	year := 0

	fmt.Println("Type a random year: ")
	fmt.Scanf("%d", &year)

	leapYear(year)
}
