package main

import "fmt"

func main() {

	var day int = 0

	fmt.Println("Diz que dia da semana é hoje?")
	fmt.Scanf("%d", &day)

	switch {
	case day == 1:
		fmt.Println("Sunday")
	case day == 2:
		fmt.Println("Monday")
	case day == 3:
		fmt.Println("Tuesday")
	case day == 4:
		fmt.Println("Wednesday")
	case day == 5:
		fmt.Println("Thursday")
	case day == 6:
		fmt.Println("Friday")
	case day == 7:
		fmt.Println("Saturday")
	default:
		fmt.Println("Este dia da Semana não existe.")
	}
}
