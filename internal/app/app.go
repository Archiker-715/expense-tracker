package app

import (
	"os"

	"github.com/Archiker-715/expense-tracker/internal/constants"
	exp "github.com/Archiker-715/expense-tracker/internal/expense"
)

func Run() {

	defer exp.CheckBudget()

	// for dbg

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"add",
	// 	"--description", "desc",
	// 	"--amount", "100",
	// 	"--test1", "100",
	// 	"--category", "another",
	// 	"--test2", "test",
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
	// 	"--id",
	// 	"--deSCRiption",
	// 	"--test1",
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
	// 	"--month", "11",
	// 	"--year", "2025",
	// 	"--amount",
	// 	"--test1",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"setbudget",
	// 	"--month", "11",
	// 	"--budget", "1",
	// 	"--checkcol", "t1",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"updatebudget",
	// 	"--month", "12",
	// 	"--budget", "1111",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"listbudget",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"deletebudget",
	// 	"--month", "12",
	// }

	// os.Args = []string{
	// 	"C:\\Users\\629B~1\\AppData\\Local\\Temp\\go-build3287855105\\b001\\exe\\main.exe",
	// 	"export",
	// 	"--id",
	// 	"--deSCRiption",
	// 	"--test1",
	// 	"--test2",
	// 	"--category",
	// }

	os.Args = os.Args[1:]

	switch os.Args[0] {
	case constants.Add:
		addExp(os.Args)
	case constants.Update:
		updateExp(os.Args)
	case constants.Delete:
		deleteExp(os.Args)
	case constants.DeleteCategory:
		deleteCategory(os.Args)
	case constants.List:
		list(os.Args)
	case constants.Summary:
		summary(os.Args)
	case constants.SetBudget:
		setBudget(os.Args)
	case constants.UpdateBudget:
		updateBudget(os.Args)
	case constants.ListBudget:
		listBudget(os.Args)
	case constants.DeleteBudget:
		deleteBudget(os.Args)
	case constants.Export:
		exportCSV(os.Args)
	}
}
