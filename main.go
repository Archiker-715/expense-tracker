package main

import (
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
	// 	"--test1", "100",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"update",
	// 	"--id", "2",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"delete",
	// 	"--id", "8",
	// }

	var (
		flags []string
		err   error
	)

	switch os.Args[1] {
	case constants.Add:
		if len(os.Args) > 2 {
			if flags, err = parse(os.Args); err != nil {
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
			if flags, err = parse(os.Args); err != nil {
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
			if flags, err = parse(os.Args); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("empty flags list")
		}

		if err := exp.DeleteExpense(flags); err != nil {
			log.Fatal(err)
		}
	}
}

func parse(userInput []string) (flags []string, err error) {
	duplicateFlags := func(flags []string) error {
		for i, flag := range flags {
			if i%2 == 0 {
				idx := slices.Index(flags, flag)
				secondIdx := slices.Index(flags[idx+1:], flag)
				if secondIdx != -1 {
					return fmt.Errorf("duplicate flag %q", flag)
				}
			}
		}
		return nil
	}

	flags = make([]string, 0)
	for i, str := range userInput[2:] {
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

	if err := duplicateFlags(flags); err != nil {
		return nil, fmt.Errorf("duplicate check: %w", err)
	}

	return
}
