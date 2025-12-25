package exp

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Archiker-715/expense-tracker/internal/constants"
	"github.com/Archiker-715/expense-tracker/internal/entity"
	fm "github.com/Archiker-715/expense-tracker/internal/file-manager"
)

func AddOpt(flags []string) error {
	sortOpts := func(opts entity.Opts) entity.Opts {
		var sortedOpts entity.Opts
		for i := 1; i < 13; i++ {
			for _, budget := range opts.Budget {
				if i == budget.Month {
					sortedOpts.Budget = append(sortedOpts.Budget, budget)
				}
			}
		}
		return sortedOpts
	}

	file, budget, opts, err := fm.PrepareJSON(flags, constants.SetBudget)
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
	file, budget, opts, err := fm.PrepareJSON(flags, constants.UpdateBudget)
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
	_, _, opts, err := fm.PrepareJSON(nil, "")
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
	file, budget, opts, err := fm.PrepareJSON(flags, constants.DeleteBudget)
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

func CheckBudget() error {
	jsonFile, err := fm.Open(constants.OptionsFileName, os.O_RDONLY)
	if err != nil {
		return fmt.Errorf("open json: %w", err)
	}
	defer jsonFile.Close()

	parsedTime, err := time.Parse(time.DateTime, time.Now().Format(time.DateTime))
	if err != nil {
		return fmt.Errorf("parse time: %w", err)
	}

	year := parsedTime.Year()
	month := int(parsedTime.Month())

	filter := map[string]string{
		constants.Month: strconv.Itoa(month),
		constants.Year:  strconv.Itoa(year),
	}

	var opts entity.Opts
	if err := json.Unmarshal(fm.ReadJson(jsonFile), &opts); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}
	if len(opts.Budget) == 0 {
		return nil
	}

	var (
		checkColumn string
		budgetSum   int
	)
	for _, budget := range opts.Budget {
		if budget.Month == month {
			checkColumn = budget.ColumnCheck
			budgetSum = budget.BudgetSum
			break
		}
	}
	if checkColumn == "" {
		return nil
	}
	checkColumn = strings.ToUpper(checkColumn[:1]) + strings.ToLower(checkColumn[1:])
	summaryFlags := []string{checkColumn}

	err, flagData := Summary(summaryFlags, filter)
	if err != nil {
		return fmt.Errorf("summary: %w", err)
	}

	for _, v := range flagData {
		if v.Sum > budgetSum {
			fmt.Printf("Warning: for column %q exceeded budget limit. Expenses: %d , budget: %d, difference: %d ", v.Flag, v.Sum, budgetSum, budgetSum-v.Sum)
		}
	}

	return nil
}
