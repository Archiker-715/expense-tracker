package exp

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Archiker-715/expense-tracker/constants"
	fm "github.com/Archiker-715/expense-tracker/file-manager"
)

type Budget struct {
	BudgetSum   int    `json:"budgetSum"`
	Month       int    `json:"month"`
	ColumnCheck string `json:"columnCheck"`
}

type Opts struct {
	Budget []Budget `json:"budget"`
}

func AddOpt(flags []string) error {
	sortOpts := func(opts Opts) Opts {
		var sortedOpts Opts
		for i := 1; i < 13; i++ {
			for _, budget := range opts.Budget {
				if i == budget.Month {
					sortedOpts.Budget = append(sortedOpts.Budget, budget)
				}
			}
		}
		return sortedOpts
	}

	file, budget, opts, err := prepareJSON(flags, constants.SetBudget)
	if err != nil {
		return fmt.Errorf("prepare JSON: %w", err)
	}
	defer file.Close()

	opts.Budget = append(opts.Budget, budget)
	sortedOpts := sortOpts(opts)

	b, err := json.MarshalIndent(sortedOpts, "", " ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}
	if err := fm.Write(file, os.O_RDWR, b); err != nil {
		return fmt.Errorf("create budget: %w", err)
	}

	return nil
}

func UpdateOpt(flags []string) error {
	file, budget, opts, err := prepareJSON(flags, constants.UpdateBudget)
	if err != nil {
		return fmt.Errorf("prepare JSON: %w", err)
	}
	defer file.Close()

	for i := 0; i < len(opts.Budget); i++ {
		if opts.Budget[i].Month == budget.Month {
			if budget.BudgetSum != 0 {
				opts.Budget[i].BudgetSum = budget.BudgetSum
			}
			if budget.ColumnCheck != "" {
				opts.Budget[i].ColumnCheck = budget.ColumnCheck
			}
		}
	}

	b, err := json.MarshalIndent(opts, "", " ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}
	if err := fm.Write(file, os.O_RDWR, b); err != nil {
		return fmt.Errorf("create budget: %w", err)
	}

	return nil
}

func ListOpt() error {
	_, _, opts, err := prepareJSON(nil, "")
	if err != nil {
		return fmt.Errorf("prepare JSON: %w", err)
	}
	b, err := json.MarshalIndent(opts, "", " ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	fmt.Println(string(b))

	return nil
}

func DeleteOpt(flags []string) error {
	file, budget, opts, err := prepareJSON(flags, constants.DeleteBudget)
	if err != nil {
		return fmt.Errorf("prepare JSON: %w", err)
	}
	defer file.Close()

	for i, b := range opts.Budget {
		if b.Month == budget.Month {
			opts.Budget = slices.Delete(opts.Budget, i, i+1)
		}
	}

	b, err := json.MarshalIndent(opts, "", " ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}
	if err := fm.Write(file, os.O_RDWR, b); err != nil {
		return fmt.Errorf("create budget: %w", err)
	}

	return nil
}

func prepareJSON(flags []string, operation string) (file *os.File, budget Budget, opts Opts, err error) {
	checkMonthExist := func(opts Opts, month int) bool {
		for _, budget := range opts.Budget {
			if budget.Month == month {
				return true
			}
		}
		return false
	}

	exists := fm.CheckExist(constants.OptionsFileName)
	if exists {
		if file, err = fm.Open(constants.OptionsFileName, os.O_RDWR); err != nil {
			return nil, budget, opts, fmt.Errorf("create %q: %w", constants.OptionsFileName, err)
		}
		if err = json.Unmarshal(fm.ReadJson(file), &opts); err != nil {
			return nil, budget, opts, fmt.Errorf("unmarshall err: %w", err)
		}
	} else {
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.OptionsFileName)
		if file, err = fm.Create(constants.OptionsFileName); err != nil {
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
