package models

// This represents the summary information for a set of transaction
type Summary struct {
	SummaryID               int
	AccountID               int
	TotalBalance            float64
	TotalTransactions       int
	NumOfCreditTransactions int
	NumOfDebitTransactions  int
	TotalAverageCredit      float64
	TotalAverageDebit       float64
}

// This represents month summary
type MonthSummary struct {
	MonthSummaryID          int
	Month                   string
	TotalBalance            float64
	TotalTransactions       int
	NumOfCreditTransactions int
	NumOfDebitTransactions  int
	AverageCredit           float64
	AverageDebit            float64
	SummaryID               int
}
