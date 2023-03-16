package models

// This represents the summary information for a set of transaction
type Summary struct {
	ID                      int
	TotalBalance            float64
	TotalTransactions       int
	NumOfCreditTransactions int
	NumOfDebitTransactions  int
	TotalAverageCredit      float64
	TotalAverageDebit       float64
	TransactionID           int
}

// This represents month summary
type MonthSummary struct {
	ID                      int
	Month                   string
	TotalBalance            float64
	TotalTransactions       int
	NumOfCreditTransactions int
	NumOfDebitTransactions  int
	AverageCredit           float64
	AverageDebit            float64
	SummaryId               int
}
