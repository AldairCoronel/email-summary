package models

// This represents the summary information for a set of transaction
type Summary struct {
	Id                     int
	TotalBalance           float64
	NumOfCreditTansactions int
	NumOfDebitTansactions  int
	TotalAverageCredit     float64
	TotalAverageDebit      float64
}

// This represents month summary
type MonthSummary struct {
	Id                     int
	Month                  string
	NumOfCreditTansactions int
	NumOfDebitTansactions  int
	AverageCredit          float64
	AverageDebit           float64
	SummaryId              int
}
