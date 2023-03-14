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

	// MonthSummaryRepository methods
	SaveMonthSummary(ctx context.Context, ms *models.MonthSummary) error
	GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error)
	ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error)

	// SummaryRepository methods
	SaveSummary(ctx context.Context, s *models.Summary) error
	GetSummaryByID(ctx context.Context, id int) (*models.Summary, error)

	Close() error
}

// Define the repository struct
var implementation Repository

// Create a new repository instance
func SetRepository(repository Repository) {
	implementation = repository
}

// Implement the SaveTransaction method of the Repository interface
func SaveTransaction(ctx context.Context, transaction *models.Transaction) error {
	return implementation.SaveTransaction(ctx, transaction)
}

// Implement the GetTransactionByID method of the Repository interface
func GetTransactionByID(ctx context.Context, id int) (*models.Transaction, error) {
	return implementation.GetTransactionByID(ctx, id)
}

// Implement the ListTransactions method of the Repository interface
func ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return implementation.ListTransactions(ctx)
}

// Implement the SaveMonthSummary method of the Repository interface
func SaveMonthSummary(ctx context.Context, ms *models.MonthSummary) error {
	return implementation.SaveMonthSummary(ctx, ms)
}

// Implement the GetMonthSummaryByID method of the Repository interface
func GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error) {
	return implementation.GetMonthSummaryByID(ctx, id)
}

// Implement the ListMonthSummaries method of the Repository interface
func ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error) {
	return implementation.ListMonthSummaries(ctx)
}

// Implement the SaveSummary method of the Repository interface
func SaveSummary(ctx context.Context, s *models.Summary) error {
	return implementation.SaveSummary(ctx, s)
}

// Implement the GetSummaryByID method of the Repository interface
func GetSummaryByID(ctx context.Context, id int) (*models.Summary, error) {
	return implementation.GetSummaryByID(ctx, id)
}

// Implement the Close method of the Repository interface
func Close() error {
	return implementation.Close()
}
