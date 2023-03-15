package repository

import (
	"context"

	"github.com/aldaircoronel/email-summary/internal/models"
)

// Define the Repository interface
type Repository interface {

	// TransactionRepository methods
	SaveTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByID(ctx context.Context, id int) (*models.Transaction, error)
	ListTransactions(ctx context.Context) ([]*models.Transaction, error)

	// SummaryRepository methods
	SaveSummary(ctx context.Context, s *models.Summary) error
	GetSummaryByID(ctx context.Context, id int) (*models.Summary, error)

	// MonthSummaryRepository methods
	SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error
	GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error)
	GetMonthSummaryBySummaryID(ctx context.Context, summaryID int) ([]*models.MonthSummary, error)
	ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error)

	Close() error
}

// Define the repository struct
var implementation Repository

// SetRepository sets the global repository implementation
func SetRepository(repository Repository) {
	implementation = repository
}

// SaveTransaction saves the given transaction
func SaveTransaction(ctx context.Context, transaction *models.Transaction) error {
	return implementation.SaveTransaction(ctx, transaction)
}

// GetTransactionByID retrieves the transaction with the given ID
func GetTransactionByID(ctx context.Context, id int) (*models.Transaction, error) {
	return implementation.GetTransactionByID(ctx, id)
}

// ListTransactions retrieves a list of all transactions
func ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return implementation.ListTransactions(ctx)
}

// SaveSummary saves the given summary
func SaveSummary(ctx context.Context, s *models.Summary) error {
	return implementation.SaveSummary(ctx, s)
}

// GetSummaryByID retrieves the summary with the given ID
func GetSummaryByID(ctx context.Context, id int) (*models.Summary, error) {
	return implementation.GetSummaryByID(ctx, id)
}

// SaveMonthSummary saves the given month summary for the given summary ID
func SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error {
	return implementation.SaveMonthSummary(ctx, ms, summaryID)
}

// GetMonthSummaryByID retrieves the month summary with the given ID
func GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error) {
	return implementation.GetMonthSummaryByID(ctx, id)
}

// GetMonthSummaryBySummaryID retrieves a list of month summaries for the given summary ID
func GetMonthSummaryBySummaryID(ctx context.Context, summaryID int) ([]*models.MonthSummary, error) {
	return implementation.GetMonthSummaryBySummaryID(ctx, summaryID)
}

// Implement the ListMonthSummaries method of the Repository interface
func ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error) {
	return implementation.ListMonthSummaries(ctx)
}

// Implement the Close method of the Repository interface
func Close() error {
	return implementation.Close()
}
