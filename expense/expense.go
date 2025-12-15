package exp

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Archiker-715/expense-tracker/constants"
	fm "github.com/Archiker-715/expense-tracker/file-manager"
)

func AddExpense(expenseDesc, expenseAmount *string, untypedFlags []string) (err error) {
	maxExpId := func(slice [][]string) (string, error) {
		if len(slice) == 1 {
			return "1", nil // len == 1 because csv have only headers
		}
		s := slice[len(slice)-1]
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

	initHeaders := func(untypedFlags []string) (headers [][]string, values []string) {
		initialHeaders := []string{"ID", "Date", "Description", "Amount"}
		headers = make([][]string, 0, (len(untypedFlags)/2)+len(initialHeaders))
		values = make([]string, 0, (len(untypedFlags) / 2))
		for i, v := range untypedFlags {
			if i%2 == 0 {
				initialHeaders = append(initialHeaders, v)
			} else {
				values = append(values, v)
			}
		}
		headers = append(headers, initialHeaders)
		return
	}

	fillInput := func(additionalValues []string, maxExpenseId string) [][]string {
		defaultInput := make([]string, 0)
		inp := make([][]string, 0)
		defaultInput = append(defaultInput, maxExpenseId, time.Now().Format(time.DateTime), *expenseDesc, *expenseAmount)
		defaultInput = append(defaultInput, additionalValues...)
		inp = append(inp, defaultInput)
		return inp
	}

	var (
		file             *os.File
		CSVheaders       [][]string
		additionalValues []string
	)
	CSVheaders, additionalValues = initHeaders(untypedFlags)
	switch fm.CheckExist(constants.ExpenseFileName) {
	case false:
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.ExpenseFileName)
		if file, err = fm.Create(constants.ExpenseFileName); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		fm.Write(file, CSVheaders)
		fmt.Printf("file %q succesfully created\n", constants.ExpenseFileName)
		file.Close()
		fallthrough
	case true:
		if file, err = fm.Open(constants.ExpenseFileName, os.O_APPEND); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
	}

	s, err := fm.Read(file)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}
	maxExpenseId, err := maxExpId(s)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}

	if len(s[0]) == len(CSVheaders[0]) {
		input := fillInput(additionalValues, maxExpenseId)
		fm.Write(file, input)
	}

	if len(s[0]) < len(CSVheaders[0]) {
		for i := 0; i < (len(CSVheaders[0]) - len(s[0])); i++ {
			for j := 0; j < len(s); j++ {
				s[j] = append(s[j], "")
			}
		}
	}

	return nil
}
