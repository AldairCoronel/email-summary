package models

// This represents month summary
type MonthSummary struct {
	Year              int // optional
	Month             string
	NumOfTransactions int
	AverageCredit     float64
	AverageDebit      float64
}

// This represents the summary information for a set of transaction
type Summary struct {
	TotalBalance           float64
	NumOfTotalTransactions int
	MonthlySummaries       []MonthSummary
}
