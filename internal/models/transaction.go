package models

import "time"

// This represents a transaction
type Transaction struct {
	TransactionID int
	AccountID     int
	ID            int
	Date          time.Time
	Amount        float64
	IsCredit      bool
}
