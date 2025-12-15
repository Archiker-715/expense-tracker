package exp

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Archiker-715/expense-tracker/constants"
	fm "github.com/Archiker-715/expense-tracker/file-manager"
)

func AddExpense(expenseDesc, expenseAmount *string, untypedFlags []string) (err error) {
	maxExpId := func(s []string) (string, error) {
		var maxExpenseId int
		for range s {
			v, err := strconv.Atoi(s[0])
			if err != nil {
				return "0", fmt.Errorf("getting maxId: %w", err)
			}
			if v > maxExpenseId {
				maxExpenseId = v
				break
			}
		}
		return strconv.Itoa(maxExpenseId + 1), nil
	}

	var file *os.File
	if !fm.CheckExist(constants.ExpenseFileName) {
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.ExpenseFileName)
		if file, err = fm.Create(constants.ExpenseFileName); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		fmt.Printf("file %q succesfully created\n", constants.ExpenseFileName)
	} else {
		if file, err = fm.Open(constants.ExpenseFileName, os.O_APPEND); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
	}

	s, err := fm.Read(file)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}
	maxExpenseId, err := maxExpId(s[len(s)-1])
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}
	input := [][]string{
		{maxExpenseId, *expenseDesc, *expenseAmount},
	}

	fm.Write(file, input)

	return nil
}
