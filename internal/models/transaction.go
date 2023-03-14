package models

import "time"

// This represents a debit or credit card transaction on an account
type Transaction struct {
	Id       int
	Date     time.Time
	Amount   float64
	IsCredit bool
}
