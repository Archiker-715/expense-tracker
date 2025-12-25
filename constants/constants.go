package constants

const (
	ExpenseFileName         = "expenses.csv"
	ExportedExpenseFileName = "export-expenses.csv"
	OptionsFileName         = "expense-options.json"

	// commands
	Add            = "add"
	Update         = "update"
	Delete         = "delete"
	DeleteCategory = "delcat"
	List           = "list"
	Summary        = "summary"
	SetBudget      = "setbudget"
	UpdateBudget   = "updatebudget"
	ListBudget     = "listbudget"
	DeleteBudget   = "deletebudget"
	Export         = "export"

	// flags sysnames
	Budget = "budget"
	Columm = "checkcol"
	Month  = "month"
	Year   = "year"

	// typed columns
	Id   = "ID"
	Date = "Date"
)
