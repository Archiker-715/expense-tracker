package app

import (
	"fmt"
	"log"

	exp "github.com/Archiker-715/expense-tracker/internal/expense"
	fp "github.com/Archiker-715/expense-tracker/internal/flags"
)

var (
	flags []string
	err   error
)

func addExp(args []string) {
	args = args[1:]
	if len(args) >= 2 { // minimum is "--column", "value"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}

	if err = exp.AddExpense(flags); err != nil {
		log.Fatal(err)
	}
}

func updateExp(args []string) {
	args = args[1:]
	if len(args) >= 4 { // minimum is "--id", "0", "--column", "value"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}

	if err := exp.UpdateExpense(flags); err != nil {
		log.Fatal(err)
	}
}

func deleteExp(args []string) {
	args = args[1:]
	if len(args) == 2 { // must be "--id", "0"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}

	if err := exp.DeleteExpense(flags); err != nil {
		log.Fatal(err)
	}
}

func list(args []string) {
	if len(args) == 1 { // only list
		if _, err := exp.ListExpense(nil); err != nil {
			log.Fatal(err)
		}
	} else if len(args) > 1 { // list + flags: "--id", "--etc"
		if flags, err = fp.Parse(args, true); err != nil {
			log.Fatal(err)
		}
		if _, err := exp.ListExpense(flags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}
}

func deleteCategory(args []string) {
	args = args[1:]
	if len(args) >= 1 { // minimum is "--column", "value"
		if flags, err = fp.Parse(args, true); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}

	if err := exp.DeleteCategories(flags); err != nil {
		log.Fatal(err)
	}
}

func summary(args []string) {
	args = args[1:]
	if len(args) >= 1 { // minimum is "--column"
		var dateFilter map[string]string
		args, dateFilter = fp.DateFilters(args)
		if len(args) >= 1 { // DateFilters could delete --month and --year if exists, then minimum is "--column"
			if flags, err = fp.Parse(args, true); err != nil {
				log.Fatal(err)
			}
			err, flagData := exp.Summary(flags, dateFilter)
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range flagData {
				fmt.Printf("Columm %q, Summary: %d\n", v.Flag, v.Sum)
			}
		} else {
			log.Fatalf("empty flags list")
		}
	} else {
		log.Fatalf("empty flags list")
	}
}

func setBudget(args []string) {
	args = args[1:]
	if len(args) == 6 { // must be "--month", "11", "--budget", "1", "--checkcol", "column"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
		if err := exp.AddOpt(flags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("not enough flags: need budget, month, checkcol")
	}
}

func updateBudget(args []string) {
	args = args[1:]
	if len(args) == 4 { // must be "--month", "12", --budget", "1111"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
		if err := exp.UpdateOpt(flags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}

}

func listBudget(args []string) {
	args = args[1:]
	if len(args) == 1 { // only listbudget
		if err := exp.ListOpt(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}
}

func deleteBudget(args []string) {
	args = args[1:]
	if len(args) == 2 { // must be "--month", "12"
		if flags, err = fp.Parse(args, false); err != nil {
			log.Fatal(err)
		}
		if err := exp.DeleteOpt(flags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}
}

func exportCSV(args []string) {
	args = args[1:]
	if len(args) >= 1 { // minimum is "--column", "value"
		if flags, err = fp.Parse(args, true); err != nil {
			log.Fatal(err)
		}
		if err := exp.Export(flags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("empty flags list")
	}
}
