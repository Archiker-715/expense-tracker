package main

import (
	"flag"
	"io"
	"log"
	"os"
	"slices"
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
		"--test2", "20",
	}

	switch os.Args[1] {
	case constants.Add:
		add := flag.NewFlagSet("add", flag.ContinueOnError)
		addDescription := add.String(constants.Description, "", "expense description")
		addAmount := add.String(constants.Amount, "0", "expense amount")

		add.SetOutput(io.Discard)
		var untypedFlags []string
		if err := add.Parse(os.Args[2:]); err != nil {
			if len(os.Args) > 0 {
				untypedFlags = parseUntypedFlags(os.Args, *addDescription, *addAmount, constants.Add)
			} else {
				log.Fatalf("parsing flags: %v", err)
			}
		}

		if err := exp.AddExpense(addDescription, addAmount, untypedFlags); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func parseUntypedFlags(s []string, expenseDesc, expenseAmount, command string) (untypedFlags []string) {
	output := make([]string, 0)
	for _, str := range s[1:] {
		if !strings.Contains(str, command) && !strings.Contains(str, constants.Description) && !strings.Contains(str, expenseDesc) && !strings.Contains(str, constants.Amount) {
			if strings.Contains(str, "--") {
				str = strings.TrimLeft(str, "-")
				output = append(output, str)
				continue
			}
			output = append(output, str)
		}
	}
	idx := slices.Index(output, expenseAmount)
	untypedFlags = slices.Delete(output, idx, idx+1)
	if len(untypedFlags)%2 != 0 {
		log.Fatalf("test")
	}

	return
}
