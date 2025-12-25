package entity

type Budget struct {
	BudgetSum   int    `json:"budgetSum"`
	Month       int    `json:"month"`
	ColumnCheck string `json:"columnCheck"`
}

type Opts struct {
	Budget []Budget `json:"budget"`
}
