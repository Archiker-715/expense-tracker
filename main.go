package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Archiker-715/expense-tracker/constants"
	exp "github.com/Archiker-715/expense-tracker/expense"
	fm "github.com/Archiker-715/expense-tracker/file-manager"
)

func main() {

	// for dbg

	os.Args = []string{
		"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
		"add",
		"--description", "desc",
		"--amount", "100",
		"--test1", "100",
	}

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"update",
	// 	"--id", "2",
	// 	"--description", "description",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"delete",
	// 	"--id", "2",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"list",
	// 	"id",
	// 	"deSCRiption",
	// 	"test1",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"delcat",
	// 	"--test3",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"summary",
	// 	"--test1",
	// 	"--test2",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"summary",
	// 	"--month", "12",
	// 	"--year", "2025",
	// 	"--amount",
	// 	"--test1",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"setbudget",
	// 	"--month", "12",
	// 	"--budget", "1",
	// 	"--checkcol", "amount",
	// }

	// check out about export csv with filters
	// todo: CRUD for json

	var (
		flags []string
		err   error
	)

	defer checkBudget()

	os.Args = os.Args[1:]

	switch os.Args[0] {
	case constants.Add:
		if len(os.Args) > 2 {
			if flags, err = parse(os.Args, false); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}

		if err = exp.AddExpense(flags); err != nil {
			log.Fatal(err)
		}
	case constants.Update:
		if len(os.Args) > 2 {
			if flags, err = parse(os.Args, false); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}

		if err := exp.UpdateExpense(flags); err != nil {
			log.Fatal(err)
		}
	case constants.Delete:
		if len(os.Args) > 2 {
			if flags, err = parse(os.Args, false); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}

		if err := exp.DeleteExpense(flags); err != nil {
			log.Fatal(err)
		}
	case constants.DeleteCategory:
		if len(os.Args) > 1 {
			if flags, err = parse(os.Args, true); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}

		if err := exp.DeleteCategories(flags); err != nil {
			log.Fatal(err)
		}
	case constants.List:
		if len(os.Args) == 1 {
			if err := exp.ListExpense(nil); err != nil {
				log.Fatal(err)
			}
		} else if len(os.Args) > 1 {
			if flags, err = parse(os.Args, true); err != nil {
				log.Fatal(err)
			}
			if err := exp.ListExpense(flags); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}
	case constants.Summary:
		if len(os.Args) > 1 {
			var dateFilter map[string]string
			os.Args, dateFilter = dateFilters(os.Args)
			if len(os.Args) > 1 {
				if flags, err = parse(os.Args, true); err != nil {
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
	case constants.SetBudget:
		if len(os.Args) == 7 {
			if flags, err = parse(os.Args, false); err != nil {
				log.Fatal(err)
			}
			if err := exp.СreateOpts(flags); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("not enough flags: need budget, month, checkcol")
		}
	}
}

func parse(userInput []string, haveOnlyFlags bool) (flags []string, err error) {
	duplicateFlags := func(flags []string) bool {
		for i, flag := range flags {
			if i%2 == 0 {
				idx := slices.Index(flags, flag)
				secondIdx := slices.Index(flags[idx+1:], flag)
				if secondIdx != -1 {
					return false
				}
			}
		}
		return true
	}

	flags = make([]string, 0)
	if !haveOnlyFlags {
		for i, str := range userInput[1:] {
			if i%2 == 0 {
				if strings.Contains(str, "--") {
					str = strings.TrimLeft(str, "-")
					if strings.EqualFold(strings.ToUpper(str), strings.ToUpper(constants.Id)) {
						str = strings.ToUpper(str)
						flags = append(flags, str)
						continue
					}
					str = strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
					flags = append(flags, str)
					continue
				} else {
					return nil, fmt.Errorf("parsing flags error on value %q", str)
				}
			} else if i%2 != 0 {
				if !strings.Contains(str, "--") {
					flags = append(flags, str)
					continue
				} else {
					return nil, fmt.Errorf("parsing flags error on value %q", str)
				}
			}
		}
		if len(flags)%2 != 0 {
			return nil, fmt.Errorf("pair flags and value error. Your input %q, parsing result %q", userInput, flags)
		}
	} else {
		for _, str := range userInput[1:] {
			if strings.Contains(str, "--") {
				str = strings.TrimLeft(str, "-")
				if userInput[0] == constants.DeleteCategory {
					if strings.EqualFold(strings.ToUpper(str), strings.ToUpper(constants.Id)) || strings.EqualFold(strings.ToUpper(str), strings.ToUpper(constants.Date)) {
						return nil, fmt.Errorf("cannot delete %q column", str)
					}
				}
				str = strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
				flags = append(flags, str)
			} else {
				return nil, fmt.Errorf("parsing flags error on value %q", str)
			}
		}
	}

	if double := duplicateFlags(flags); !double {
		return nil, errors.New("duplicate check: input have double of flag")
	}

	return
}

func dateFilters(userInput []string) ([]string, map[string]string) {
	date := make(map[string]string)
	idxs := make([]int, 0)
	for i, str := range userInput {
		if strings.Contains(str, "--") {
			str = strings.TrimLeft(str, "-")
			if strings.EqualFold(str, constants.Month) || strings.EqualFold(str, constants.Year) {
				if i+1 <= len(userInput)-1 {
					date[str] = userInput[i+1]
					idxs = append(idxs, i, i+1)
				}
			}
		}
	}

	revertIdxs := make([]int, 0, len(idxs))
	for i := len(idxs) - 1; i >= 0; i-- {
		revertIdxs = append(revertIdxs, idxs[i])
	}

	for _, revertIdx := range revertIdxs {
		userInput = slices.Delete(userInput, revertIdx, revertIdx+1)
	}

	return userInput, date

}

func checkBudget() error {
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

	var opts exp.Opts
	if err := json.Unmarshal(fm.ReadJson(jsonFile), &opts); err != nil {
		return fmt.Errorf("parse json: %w", err)
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
	checkColumn = strings.ToUpper(checkColumn[:1]) + strings.ToLower(checkColumn[1:])
	summaryFlags := []string{checkColumn}

	err, flagData := exp.Summary(summaryFlags, filter)
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
