package exp

import (
	"encoding/json"
	"fmt"
	"os"
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

func СreateOpts(flags []string) error {

	checkMonthExist := func(opts Opts, month int) bool {
		for _, budget := range opts.Budget {
			if budget.Month == month {
				return true
			}
		}
		return false
	}

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

	var file *os.File
	var err error
	switch fm.CheckExist(constants.OptionsFileName) {
	case false:
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.OptionsFileName)
		if file, err = fm.Create(constants.OptionsFileName); err != nil {
			return fmt.Errorf("create %q: %w", constants.OptionsFileName, err)
		}
		fmt.Printf("file %q succesfully created\n", constants.OptionsFileName)
		file.Close()
		fallthrough
	case true:
		if file, err = fm.Open(constants.OptionsFileName, os.O_RDWR); err != nil {
			return fmt.Errorf("create %q: %w", constants.OptionsFileName, err)
		}
	}
	defer file.Close()

	var (
		opts   Opts
		budget Budget
		month  int
	)

	for i := 0; i < len(flags)-1; i++ {
		if strings.EqualFold(flags[i], constants.Budget) {
			v, err := strconv.Atoi(flags[i+1])
			if err != nil {
				return fmt.Errorf("convert budget to int: %w", err)
			}
			budget.BudgetSum = v
		}
		if strings.EqualFold(flags[i], constants.Month) {
			v, err := strconv.Atoi(flags[i+1])
			if err != nil {
				return fmt.Errorf("convert month to int: %w", err)
			}
			if v <= 0 || v >= 13 {
				return fmt.Errorf("month must be in 1-12. Your input: '%d'", v)
			}
			budget.Month, month = v, v
		}
		if strings.EqualFold(flags[i], constants.Columm) {
			v := flags[i+1]
			v = strings.ToUpper(v[:1]) + strings.ToLower(v[1:])
			budget.ColumnCheck = v
		}
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("fileStat err: %w", err)
	}

	if int(fileInfo.Size()) != 0 {
		if err = json.Unmarshal(fm.ReadJson(file), &opts); err != nil {
			return fmt.Errorf("unmarshall err: %w", err)
		}
	}
	if exists := checkMonthExist(opts, month); exists {
		return fmt.Errorf("month '%d' already exists", month)
	}

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
