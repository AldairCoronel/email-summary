package models

// This represents month summary
type MonthSummary struct {
	Month             string
	NumOfTransactions int
	AverageCredit     float64
	AverageDebit      float64
}

type TransactionByMonth struct {
	Month             string
	NumOfTransactions int
}

// This represents the summary information for a set of transaction
type Summary struct {
	TotalBalance           float64
	NumOfTotalTransactions int
	TransactionByMonth     []TransactionByMonth
	MonthlySummaries       []MonthSummary
}
