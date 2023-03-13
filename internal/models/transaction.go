package models

import "time"

// This represents a debit or credit card transaction on an account
type Transaction struct {
	ID        int
	Date      time.Time
	Amount    float64
	IsCredit  bool
	MonthYear string
}
