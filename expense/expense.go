package exp

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Archiker-715/expense-tracker/constants"
	fm "github.com/Archiker-715/expense-tracker/file-manager"
)

func AddExpense(flags []string) (err error) {
	maxExpId := func(slice [][]string) (string, error) {
		if len(slice) == 1 {
			return "1", nil // len == 1 because csv have only headers
		}
		s := slice[len(slice)-1]
		var maxExpenseId int
		for range s {
			v, err := strconv.Atoi(s[0])
			if err != nil {
				return "0", fmt.Errorf("getting maxId: %w", err)
			}
			if v > maxExpenseId {
				maxExpenseId = v
				break
			}
		}
		return strconv.Itoa(maxExpenseId + 1), nil
	}

	initHeaders := func(untypedFlags []string) (headers [][]string, values []string) {
		initialHeaders := []string{constants.Id, constants.Date}
		headers = make([][]string, 0, (len(untypedFlags)/2)+len(initialHeaders))
		values = make([]string, 0, (len(untypedFlags) / 2))
		for i, v := range untypedFlags {
			if i%2 == 0 {
				initialHeaders = append(initialHeaders, v)
			} else {
				values = append(values, v)
			}
		}
		headers = append(headers, initialHeaders)
		return
	}

	fillInput := func(additionalValues []string, maxExpenseId string) [][]string {
		defaultInput := make([]string, 0)
		inp := make([][]string, 0)
		defaultInput = append(defaultInput, maxExpenseId, time.Now().Format(time.DateTime))
		defaultInput = append(defaultInput, additionalValues...)
		inp = append(inp, defaultInput)
		return inp
	}

	addHeaderWriteInput := func(CSVheaders, inputCSVheaders, input [][]string, file *os.File) error {
		iCondition := len(inputCSVheaders[0]) - len(CSVheaders[0])
		for i := 0; i < iCondition; i++ {
			for j := 0; j < len(CSVheaders); j++ {
				if j == 0 {
					CSVheaders[j] = inputCSVheaders[0]
				} else {
					CSVheaders[j] = append(CSVheaders[j], "")
				}
			}
		}
		CSVheaders = append(CSVheaders, input[0])
		file.Close()
		if file, err = fm.Open(constants.ExpenseFileName, os.O_RDWR); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		defer file.Close()

		if err := fm.Write(file, os.O_RDWR, CSVheaders); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		return nil
	}

	newHeadersFromInput := func(currentCSV, inputCSVheaders [][]string) (origHeaders, newHeaders []string, err error) {
		origHeaders, newHeaders = make([]string, 0), make([]string, 0)
		origHeaders, newHeaders = append(origHeaders, currentCSV[0]...), append(newHeaders, inputCSVheaders[0]...)

		for _, CSVheader := range currentCSV[0] {
			for _, inputCSVheader := range inputCSVheaders[0] {
				if CSVheader == inputCSVheader {
					idx := slices.Index(origHeaders, CSVheader)
					if idx == -1 {
						return nil, nil, errors.New("columns's header not found")
					}
					origHeaders = slices.Delete(origHeaders, idx, idx+1)
					idx = slices.Index(newHeaders, CSVheader)
					if idx == -1 {
						return nil, nil, errors.New("columns's header not found")
					}
					newHeaders = slices.Delete(newHeaders, idx, idx+1)
					break
				}
			}
		}
		return
	}

	addNewHeaders := func(currentCSV, input [][]string, origHeaders, newHeaders []string, file *os.File) error {
		for _, v := range origHeaders {
			idx := slices.Index(currentCSV[0], v)
			if idx == -1 {
				newHeaders = append(newHeaders, v)
			}
			input[0] = slices.Insert(input[0], idx, "")
		}

		if len(newHeaders) > 0 {
			tempNewCSVheaders := make([][]string, 0, len(currentCSV[0]))
			tempNewCSVheaders = append(tempNewCSVheaders, append(currentCSV[0], newHeaders...))
			if err := addHeaderWriteInput(currentCSV, tempNewCSVheaders, input, file); err != nil {
				return fmt.Errorf("add header: %w", err)
			}
		}

		return nil
	}

	var (
		file             *os.File
		inputCSVheaders  [][]string
		additionalValues []string
	)
	inputCSVheaders, additionalValues = initHeaders(flags)
	switch fm.CheckExist(constants.ExpenseFileName) {
	case false:
		fmt.Printf("file %q not found. Will be create in current directory\n", constants.ExpenseFileName)
		if file, err = fm.Create(constants.ExpenseFileName); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		if err = fm.Write(file, os.O_APPEND, inputCSVheaders); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		fmt.Printf("file %q succesfully created\n", constants.ExpenseFileName)
		file.Close()
		fallthrough
	case true:
		if file, err = fm.Open(constants.ExpenseFileName, os.O_APPEND); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
	}
	defer file.Close()

	currentCSV, err := fm.Read(file)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}
	maxExpenseId, err := maxExpId(currentCSV)
	if err != nil {
		return fmt.Errorf("add expense: %w", err)
	}

	eq := slices.Equal(currentCSV[0], inputCSVheaders[0])
	input := fillInput(additionalValues, maxExpenseId)

	if len(currentCSV[0]) == len(inputCSVheaders[0]) && eq {
		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}
		return nil
	}

	if len(currentCSV[0]) < len(inputCSVheaders[0]) {
		if err := addHeaderWriteInput(currentCSV, inputCSVheaders, input, file); err != nil {
			return fmt.Errorf("add header: %w", err)
		}
		return nil
	}

	if len(currentCSV[0]) == len(inputCSVheaders[0]) && !eq {
		origHeaders, newHeaders, err := newHeadersFromInput(currentCSV, inputCSVheaders)
		if err != nil {
			return fmt.Errorf("get newHeadersFromInput error: %w", err)
		}

		if err := addNewHeaders(currentCSV, input, origHeaders, newHeaders, file); err != nil {
			return fmt.Errorf("add new headers %q: %w", constants.ExpenseFileName, err)
		}

		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}

		return nil
	}

	if len(currentCSV[0]) > len(inputCSVheaders[0]) {
		origHeaders, newHeaders, err := newHeadersFromInput(currentCSV, inputCSVheaders)
		if err != nil {
			return fmt.Errorf("get newHeadersFromInput error: %w", err)
		}

		if err := addNewHeaders(currentCSV, input, origHeaders, newHeaders, file); err != nil {
			return fmt.Errorf("add new headers %q: %w", constants.ExpenseFileName, err)
		}

		if err := fm.Write(file, os.O_APPEND, input); err != nil {
			return fmt.Errorf("create %q: %w", constants.ExpenseFileName, err)
		}

		return nil
	}

	return errors.New("unexpected end of func")
}

func UpdateExpense(flags []string) error {
	buildCSV := func(csv [][]string, flagsVals []string, stringIndex int) ([][]string, error) {
		flagsIdxVals := make(map[string]map[int]string)
		idxVals := make(map[int]string)
		var tempFlag string
		var flagIdx int
		for i, val := range flagsVals {
			if i%2 == 0 {
				if flagIdx = slices.Index(csv[0], val); flagIdx != -1 {
					tempFlag = val
					idxVals[flagIdx] = ""
					flagsIdxVals[val] = idxVals
					continue
				} else {
					return nil, fmt.Errorf("entered flag %q not fount in csv", val)
				}
			} else {
				idxVals[flagIdx] = val
				flagsIdxVals[tempFlag] = idxVals
			}
		}

		for _, m := range flagsIdxVals {
			for k, v := range m {
				csv[stringIndex][k] = v
			}
		}

		return csv, nil
	}

	csv, stringIdx, file, err := prepareCSV(flags, true)
	if err != nil {
		return fmt.Errorf("prepare CSV error: %w", err)
	}
	defer file.Close()

	csv, err = buildCSV(csv, flags, stringIdx)
	if err != nil {
		return fmt.Errorf("build updated csv error: %w", err)
	}

	if err := fm.Write(file, os.O_RDWR, csv); err != nil {
		return fmt.Errorf("update csv error: %w", err)
	}

	return nil
}

func DeleteExpense(flags []string) error {

	csv, stringIdx, file, err := prepareCSV(flags, true)
	if err != nil {
		return fmt.Errorf("prepare CSV error: %w", err)
	}
	defer file.Close()

	csv = slices.Delete(csv, stringIdx, stringIdx+1)

	if err := fm.Write(file, os.O_RDWR, csv); err != nil {
		return fmt.Errorf("update csv error: %w", err)
	}

	return nil
}

func DeleteCategories(flags []string) error {
	csv, _, file, err := prepareCSV(flags, false)
	if err != nil {
		return fmt.Errorf("prepare CSV error: %w", err)
	}
	defer file.Close()

	newCSV := csvByCategory(csv, flags, constants.DeleteCategory)

	if err := fm.Write(file, os.O_RDWR, newCSV); err != nil {
		return fmt.Errorf("update csv error: %w", err)
	}

	return nil
}

func ListExpense(flags []string) error {
	csv, _, file, err := prepareCSV(flags, false)
	if err != nil {
		return fmt.Errorf("prepare CSV error: %w", err)
	}
	defer file.Close()

	if len(flags) == 0 {
		fm.Print(csv)
	} else {
		fm.Print(csvByCategory(csv, flags, constants.List))
	}

	return nil
}

func Summary(flags []string, dateFilter map[string]string) error {
	type flagMetadata struct {
		Flag string
		Sum  int
	}

	sum := func(csv [][]string, flags []string) (flagData map[int]*flagMetadata) {
		flagData = make(map[int]*flagMetadata, 0)
		for i, column := range csv[0] {
			for _, flag := range flags {
				if strings.EqualFold(column, flag) {
					flagData[i] = &flagMetadata{Flag: flag}
				}
			}
		}

		for i, csvStr := range csv {
			for j, val := range csvStr {
				for k := range flagData {
					if i > 0 {
						if k == j {
							if val != "" {
								valInt, err := strconv.Atoi(val)
								if err != nil {
									fmt.Printf("columnn: %v, string %d: cannot convert %q to int, value was not summing. Please check your CSV-file\n", flagData[k].Flag, i, val)
								}
								flagData[k].Sum += valInt
							}
						}
					}
				}
			}
		}
		return
	}

	filter := func(csv [][]string, dateFilter map[string]string) ([][]string, error) {
		var (
			yearInt  int
			monthInt int
			err      error
		)
		newCSV := make([][]string, 0)
		year, yearOk := dateFilter[constants.Year]
		if yearOk {
			if yearInt, err = strconv.Atoi(year); err != nil {
				return nil, fmt.Errorf("cannot convert %q to int", year)
			}
		}
		month, monthOk := dateFilter[constants.Month]
		if monthOk {
			if monthInt, err = strconv.Atoi(month); err != nil {
				return nil, fmt.Errorf("cannot convert %q to int", month)
			}
		}
		newCSV = append(newCSV, csv[0])
		for _, csvStr := range csv[1:] {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", csvStr[1])
			if err != nil {
				return nil, fmt.Errorf("parse time: %w", err)
			}
			if yearOk && monthOk {
				if parsedTime.Year() == yearInt && parsedTime.Month() == time.Month(monthInt) {
					newCSV = append(newCSV, csvStr)
				}
			}
			if yearOk && !monthOk {
				if parsedTime.Year() == yearInt {
					newCSV = append(newCSV, csvStr)
				}
			}
			if !yearOk && monthOk {
				if parsedTime.Month() == time.Month(monthInt) {
					newCSV = append(newCSV, csvStr)
				}
			}

		}
		return newCSV, nil
	}

	csv, _, file, err := prepareCSV(flags, false)
	if err != nil {
		return fmt.Errorf("prepare CSV error: %w", err)
	}
	defer file.Close()

	if len(dateFilter) > 0 {
		if csv, err = filter(csv, dateFilter); err != nil {
			return fmt.Errorf("filtering error: %w", err)
		}
	}

	flagData := sum(csv, flags)
	for _, v := range flagData {
		fmt.Printf("Columm %q, Summary: %d\n", v.Flag, v.Sum)
	}

	return nil
}

// base file checks and find ID in CSV
func prepareCSV(flags []string, indexingById bool) (csv [][]string, stringIdx int, file *os.File, err error) {
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

	if exists := fm.CheckExist(constants.ExpenseFileName); !exists {
		return nil, -1, file, fmt.Errorf("file %q not exists. Please add your first expense", constants.ExpenseFileName)
	}

	file, err = fm.Open(constants.ExpenseFileName, os.O_RDWR)
	if err != nil {
		return nil, -1, file, fmt.Errorf("open file error: %w", err)
	}

	csv, err = fm.Read(file)
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

func indexingCategory(CSVcolumns, flags []string) (idxs []int) {
	idxs = make([]int, 0)
	for i, column := range CSVcolumns {
		for _, flag := range flags {
			if strings.EqualFold(column, flag) {
				idxs = append(idxs, i)
			}
		}
	}
	if len(idxs) != len(CSVcolumns) {
		fmt.Println("not all columns filtered by flags. Check your input")
	}

	return
}

func csvByCategory(csv [][]string, flags []string, command string) [][]string {
	idxs := indexingCategory(csv[0], flags)

	if strings.EqualFold(command, constants.List) {
		filteredCSV := make([][]string, 0, len(csv))
		for _, csvStr := range csv {
			filteredCSVstr := make([]string, 0)
			for _, idx := range idxs {
				filteredCSVstr = append(filteredCSVstr, csvStr[idx])
			}
			filteredCSV = append(filteredCSV, filteredCSVstr)
		}
		return filteredCSV
	}

	if strings.EqualFold(command, constants.DeleteCategory) {
		newCSV := make([][]string, 0, len(csv))
		revertIdxs := make([]int, 0, len(idxs))
		for i := len(idxs) - 1; i >= 0; i-- {
			revertIdxs = append(revertIdxs, idxs[i])
		}

		for _, csvStr := range csv {
			for _, revertIdx := range revertIdxs {
				csvStr = slices.Delete(csvStr, revertIdx, revertIdx+1)
			}
			newCSV = append(newCSV, csvStr)
		}
		return newCSV
	}

	return nil
}
