package fp

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Archiker-715/expense-tracker/internal/constants"
)

func Parse(userInput []string, haveOnlyFlags bool) (flags []string, err error) {
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
		// create pairs from input like a "--column", "value"
		for i, str := range userInput {
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
		// create list from input like "--column", "--column"
		for _, str := range userInput {
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

// find --month and --year flags, save to map and delete from userInput
func DateFilters(userInput []string) ([]string, map[string]string) {
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
