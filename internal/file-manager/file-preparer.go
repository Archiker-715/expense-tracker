package fm

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Archiker-715/expense-tracker/internal/constants"
	"github.com/Archiker-715/expense-tracker/internal/entity"
)

func PrepareJSON(flags []string, operation string) (file *os.File, budget entity.Budget, opts entity.Opts, err error) {
	checkMonthExist := func(opts entity.Opts, month int) bool {
		for _, budget := range opts.Budget {
			if budget.Month == month {
				return true
			}
		}
		return false
	}

	exists := CheckExist(constants.OptionsFileName)
	if exists {
		if file, err = Open(constants.OptionsFileName, os.O_RDWR); err != nil {
			return nil, budget, opts, fmt.Errorf("create %q: %w", constants.OptionsFileName, err)
		}
		if err = json.Unmarshal(ReadJson(file), &opts); err != nil {
			return nil, budget, opts, fmt.Errorf("unmarshall err: %w", err)
		}
	} else {
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.OptionsFileName)
		if file, err = Create(constants.OptionsFileName); err != nil {
			return nil, budget, opts, fmt.Errorf("create %q: %w", constants.OptionsFileName, err)
		}
		fmt.Printf("file %q succesfully created\n", constants.OptionsFileName)

	}

	if operation == constants.ListBudget {
		file.Close()
		return nil, budget, opts, nil
	}

	var month int
	for i := 0; i < len(flags)-1; i++ {
		if strings.EqualFold(flags[i], constants.Budget) {
			v, err := strconv.Atoi(flags[i+1])
			if err != nil {
				return nil, budget, opts, fmt.Errorf("convert budget to int: %w", err)
			}
			budget.BudgetSum = v
		}
		if strings.EqualFold(flags[i], constants.Month) {
			v, err := strconv.Atoi(flags[i+1])
			if err != nil {
				return nil, budget, opts, fmt.Errorf("convert month to int: %w", err)
			}
			if v <= 0 || v >= 13 {
				return nil, budget, opts, fmt.Errorf("month must be in 1-12. Your input: '%d'", v)
			}
			budget.Month, month = v, v
		}
		if strings.EqualFold(flags[i], constants.Columm) {
			v := flags[i+1]
			v = strings.ToUpper(v[:1]) + strings.ToLower(v[1:])
			budget.ColumnCheck = v
		}
	}

	switch operation {
	case constants.SetBudget:
		if exists := checkMonthExist(opts, month); exists {
			return nil, budget, opts, fmt.Errorf("month '%d' already exists in json", month)
		}
	case constants.UpdateBudget:
		if exists := checkMonthExist(opts, month); !exists {
			return nil, budget, opts, fmt.Errorf("nothing to update. Month '%d' not exists in json", month)
		}
	case constants.DeleteBudget:
		if exists := checkMonthExist(opts, month); !exists {
			return nil, budget, opts, fmt.Errorf("nothing to delete. Month '%d' not exists in json", month)
		}
	}

	return file, budget, opts, nil
}

// base file checks and find ID in CSV
func PrepareCSV(flags []string, indexingById bool) (csv [][]string, stringIdx int, file *os.File, err error) {
	indexById := func(csv [][]string, id string) (stringIndex int) {
		for i, csvStr := range csv {
			if csvStr[0] == id {
				stringIndex = i
				break
			}
		}
		if stringIndex == 0 {
			return -1
		}
		return
	}

	var idIdx int
	if indexingById {
		idIdx = slices.Index(flags, constants.Id)
		if idIdx == -1 {
			return nil, -1, file, fmt.Errorf("nothing to update, flags not contains id")
		}
	}

	if exists := CheckExist(constants.ExpenseFileName); !exists {
		return nil, -1, file, fmt.Errorf("file %q not exists. Please add your first expense", constants.ExpenseFileName)
	}

	file, err = Open(constants.ExpenseFileName, os.O_RDWR)
	if err != nil {
		return nil, -1, file, fmt.Errorf("open file error: %w", err)
	}

	csv, err = Read(file)
	if err != nil {
		return nil, -1, file, fmt.Errorf("read csv error: %w", err)
	}

	if indexingById {
		stringIdx = indexById(csv, flags[idIdx+1])
		if stringIdx == -1 {
			return nil, -1, file, fmt.Errorf("not found 'id %v' in csv", flags[idIdx+1])
		}
	}

	return
}
