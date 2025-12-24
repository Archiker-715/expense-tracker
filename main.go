package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Archiker-715/expense-tracker/constants"
	exp "github.com/Archiker-715/expense-tracker/expense"
)

func main() {

	// for dbg

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"add",
	// 	"--description", "desc",
	// 	"--amount", "10",
	// 	"--test1", "100",
	// 	"--test3", "100",
	// }

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

	// check out about export csv with filters

	// Allow users to set a budget for each month and show a warning when the user exceeds the budget.
	// create json file with some field included budget and compare budget with current expenses for month
	// and every time i need check json-file and do summary with filters on curruent month

	var (
		flags []string
		err   error
	)

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
				if err := exp.Summary(flags, dateFilter); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatalf("empty flags list")
			}
		} else {
			log.Fatalf("empty flags list")
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
