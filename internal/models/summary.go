package models

// This represents month summary
type MonthSummary struct {
	Month             string
	Year              int // optional
	NumOfTransactions int
	AverageCredit     float64
	AverageDebit      float64
	SummaryId         int
}

// This represents the summary information for a set of transaction
type Summary struct {
	Id                     int
	TotalBalance           float64
	NumOfTotalTransactions int
	MonthlySummaries       [12]MonthSummary
}
