package exp

import (
	"errors"
	"fmt"
	"os"
	"slices"
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

	addHeaderWriteInput := func(CSVheaders, inputCSVheaders, input [][]string, file *os.File) error {
		iCondition := len(inputCSVheaders[0]) - len(CSVheaders[0])
		for i := 0; i < iCondition; i++ {
			for j := 0; j < len(CSVheaders); j++ {
				if j == 0 {
					CSVheaders[j] = inputCSVheaders[0]
				} else {
					CSVheaders[j] = append(CSVheaders[j], "")
				}
			}
		}
		CSVheaders = append(CSVheaders, input[0])
		file.Close()
		if file, err = fm.Open(constants.ExpenseFileName, os.O_RDWR); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		defer file.Close()

		if err := fm.Write(file, os.O_RDWR, CSVheaders); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		return nil
	}

	newHeadersFromInput := func(currentCSV, inputCSVheaders [][]string) (origHeaders, newHeaders []string, err error) {
		origHeaders, newHeaders = make([]string, 0), make([]string, 0)
		origHeaders, newHeaders = append(origHeaders, currentCSV[0]...), append(newHeaders, inputCSVheaders[0]...)

		for _, CSVheader := range currentCSV[0] {
			for _, inputCSVheader := range inputCSVheaders[0] {
				if CSVheader == inputCSVheader {
					idx := slices.Index(origHeaders, CSVheader)
					if idx == -1 {
						return nil, nil, errors.New("columns's header not found")
					}
					origHeaders = slices.Delete(origHeaders, idx, idx+1)
					idx = slices.Index(newHeaders, CSVheader)
					if idx == -1 {
						return nil, nil, errors.New("columns's header not found")
					}
					newHeaders = slices.Delete(newHeaders, idx, idx+1)
					break
				}
			}
		}
		return
	}

	addNewHeaders := func(currentCSV, input [][]string, origHeaders, newHeaders []string, file *os.File) error {
		for _, v := range origHeaders {
			idx := slices.Index(currentCSV[0], v)
			if idx == -1 {
				newHeaders = append(newHeaders, v)
			}
			input[0] = slices.Insert(input[0], idx, "")
		}

		if len(newHeaders) > 0 {
			tempNewCSVheaders := make([][]string, 0, len(currentCSV[0]))
			tempNewCSVheaders = append(tempNewCSVheaders, append(currentCSV[0], newHeaders...))
			if err := addHeaderWriteInput(currentCSV, tempNewCSVheaders, input, file); err != nil {
				return fmt.Errorf("add header: %w", err)
			}
		}

		return nil
	}

	var (
		file             *os.File
		inputCSVheaders  [][]string
		additionalValues []string
	)
	inputCSVheaders, additionalValues = initHeaders(untypedFlags)
	switch fm.CheckExist(constants.ExpenseFileName) {
	case false:
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.ExpenseFileName)
		if file, err = fm.Create(constants.ExpenseFileName); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		if err = fm.Write(file, os.O_APPEND, inputCSVheaders); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		fmt.Printf("file %q succesfully created\n", constants.ExpenseFileName)
		file.Close()
		fallthrough
	case true:
		if file, err = fm.Open(constants.ExpenseFileName, os.O_APPEND); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
	}

	currentCSV, err := fm.Read(file)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}
	maxExpenseId, err := maxExpId(currentCSV)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}

	eq := slices.Equal(currentCSV[0], inputCSVheaders[0])
	input := fillInput(additionalValues, maxExpenseId)

	if len(currentCSV[0]) == len(inputCSVheaders[0]) && eq {
		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		return nil
	}

	if len(currentCSV[0]) < len(inputCSVheaders[0]) {
		if err := addHeaderWriteInput(currentCSV, inputCSVheaders, input, file); err != nil {
			return fmt.Errorf("add header: %w", err)
		}
		return nil
	}

	if len(currentCSV[0]) == len(inputCSVheaders[0]) && !eq {
		origHeaders, newHeaders, err := newHeadersFromInput(currentCSV, inputCSVheaders)
		if err != nil {
			return fmt.Errorf("get newHeadersFromInput error: %w", err)
		}

		if err := addNewHeaders(currentCSV, input, origHeaders, newHeaders, file); err != nil {
			return fmt.Errorf("add new headers %q: %w", constants.ExpenseFileName, err)
		}

		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}

		return nil
	}

	if len(currentCSV[0]) > len(inputCSVheaders[0]) {
		origHeaders, newHeaders, err := newHeadersFromInput(currentCSV, inputCSVheaders)
		if err != nil {
			return fmt.Errorf("get newHeadersFromInput error: %w", err)
		}

		if err := addNewHeaders(currentCSV, input, origHeaders, newHeaders, file); err != nil {
			return fmt.Errorf("add new headers %q: %w", constants.ExpenseFileName, err)
		}

		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}

		return nil
	}

	return errors.New("unexpected end of func")
}
