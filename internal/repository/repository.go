package repository

import (
	"context"

	"github.com/aldaircoronel/email-summary/internal/models"
)

// Define the Repository interface
type Repository interface {

	// AccountRepository methods
	SaveAccount(ctx context.Context, a *models.Account) error
	GetAccountByID(ctx context.Context, id int) (*models.Account, error)

	// TransactionRepository methods
	SaveTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByAccountID(ctx context.Context, accountID int) ([]*models.Transaction, error)
	ListTransactions(ctx context.Context) ([]*models.Transaction, error)

	// SummaryRepository methods
	SaveSummary(ctx context.Context, s *models.Summary) error
	GetSummaryByAccountID(ctx context.Context, accountID int) (*models.Summary, error)
	ListSummaries(ctx context.Context) ([]*models.Summary, error)

	// MonthSummaryRepository methods
	SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error
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

// SaveAccount saves the given account
func SaveAccount(ctx context.Context, a *models.Account) error {
	return implementation.SaveAccount(ctx, a)
}

// GetAccountByID retrieves the account with the given ID
func GetAccountByID(ctx context.Context, id int) (*models.Account, error) {
	return implementation.GetAccountByID(ctx, id)
}

// SaveTransaction saves the given transaction
func SaveTransaction(ctx context.Context, transaction *models.Transaction) error {
	return implementation.SaveTransaction(ctx, transaction)
}

// GetTransactionByAccountID retrieves the transaction with the given AccountID
func GetTransactionByAccountID(ctx context.Context, accountID int) ([]*models.Transaction, error) {
	return implementation.GetTransactionByAccountID(ctx, accountID)
}

// ListTransactions retrieves a list of all transactions
func ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return implementation.ListTransactions(ctx)
}

// SaveSummary saves the given summary
func SaveSummary(ctx context.Context, s *models.Summary) error {
	return implementation.SaveSummary(ctx, s)
}

// GetSummaryByAccountID retrieves the summary with the given ID
func GetSummaryByAccountID(ctx context.Context, accountID int) (*models.Summary, error) {
	return implementation.GetSummaryByAccountID(ctx, accountID)
}

// SaveMonthSummary saves the given month summary for the given summary ID
func SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error {
	return implementation.SaveMonthSummary(ctx, ms, summaryID)
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
