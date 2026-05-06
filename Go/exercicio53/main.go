package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	cliID    int
	cliName  string
	cliFunds float64
}

func (t *BankAccount) TransferMoney(bankId *BankAccount, quantityTransfer float64) error {
	if t.cliFunds == 0 {
		return errors.New("you cannot transfer funds to another client")
	}
	if t.cliID != bankId.cliID {
		t.cliFunds -= quantityTransfer
		bankId.cliFunds += quantityTransfer
	} else {
		return errors.New("Cannot transfer your money to yourself.....")
	}
	return nil
}

func main() {
	cliA := BankAccount{
		cliID:    1,
		cliName:  "Ricardo",
		cliFunds: 22.34,
	}
	cliB := BankAccount{
		cliID:    2,
		cliName:  "Conceição",
		cliFunds: 200,
	}

	fmt.Println(cliA.cliName, ":", cliA.cliFunds)
	fmt.Println(cliB.cliName, ":", cliB.cliFunds)

	if err := cliB.TransferMoney(&cliA, 20); err != nil {
		errors.New("you cannot transfer money")

	}
	fmt.Println("Transference COMPLETED...")
	fmt.Println(cliA.cliName, ":", cliA.cliFunds)
	fmt.Println(cliB.cliName, ":", cliB.cliFunds)
}
