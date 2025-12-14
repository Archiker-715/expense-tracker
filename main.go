package main

import (
	"flag"
	"log"

	exp "github.com/Archiker-715/expense-tracker/expense"
)

func main() {
	description := flag.String("description", "", "expense description")
	amount := flag.String("amount", "0", "expense amount")
	flag.Parse()

	if err := exp.AddExpense(description, amount); err != nil {
		log.Fatalf("%v", err)
	}
}
