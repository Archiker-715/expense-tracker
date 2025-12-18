package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Archiker-715/expense-tracker/constants"
	exp "github.com/Archiker-715/expense-tracker/expense"
)

func main() {

	// for dbg
	os.Args = []string{
		"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
		"add",
		"--description", "desc",
		"--amount", "10",
		"--test1", "100",
		"--test2", "200",
	}

	switch os.Args[1] {
	case constants.Add:
		var (
			flags []string
			err   error
		)
		if len(os.Args) > 2 {
			flags, err = parse(os.Args)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err = exp.AddExpense(flags); err != nil {
			log.Fatal(err)
		}
	case constants.Update:
		flags, err := parse(os.Args)
		if err != nil {
			log.Fatal(err)
		}

		if err := exp.UpdateExpense(flags); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func parse(s []string) (flags []string, err error) {
	flags = make([]string, 0)
	for i, str := range s[2:] {
		if i%2 == 0 {
			if strings.Contains(str, "--") {
				str = strings.TrimLeft(str, "-")
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
		return nil, fmt.Errorf("pair flags and value error. Your input %q, parsing result %q", s, flags)
	}

	return
}
