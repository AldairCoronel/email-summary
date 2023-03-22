package repository

import (
	"context"
	"errors"

	"github.com/aldaircoronel/email-summary/internal/models"
)

type InMemoryRepository struct {
	accounts           map[int]*models.Account
	transactions       map[int][]*models.Transaction
	summaries          map[int]*models.Summary
	monthSummaries     map[int][]*models.MonthSummary
	nextAccountID      int
	nextTransactionID  int
	nextSummaryID      int
	nextMonthSummaryID int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts:           make(map[int]*models.Account),
		transactions:       make(map[int][]*models.Transaction),
		summaries:          make(map[int]*models.Summary),
		monthSummaries:     make(map[int][]*models.MonthSummary),
		nextAccountID:      1,
		nextTransactionID:  1,
		nextSummaryID:      1,
		nextMonthSummaryID: 1,
	}
}

// AccountRepository methods

func (r *InMemoryRepository) SaveAccount(ctx context.Context) (int, error) {
	account := &models.Account{
		AccountID: r.nextAccountID,
	}
	r.accounts[account.AccountID] = account
	r.nextAccountID++
	return account.AccountID, nil
}

func (r *InMemoryRepository) GetAccountByID(ctx context.Context, id int) (*models.Account, error) {
	account, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (r *InMemoryRepository) Close() error {
	return nil
}

// TransactionRepository methods

func (r *InMemoryRepository) SaveTransaction(ctx context.Context, trx *models.Transaction) error {
	trx.TransactionID = r.nextTransactionID
	r.nextTransactionID++
	r.transactions[trx.AccountID] = append(r.transactions[trx.AccountID], trx)
	return nil
}

func (r *InMemoryRepository) GetTransactionByAccountID(ctx context.Context, accountID int) ([]*models.Transaction, error) {
	trxs, ok := r.transactions[accountID]
	if !ok {
		return nil, errors.New("transactions not found")
	}
	return trxs, nil
}

func (r *InMemoryRepository) ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	var allTrxs []*models.Transaction
	for _, trxs := range r.transactions {
		allTrxs = append(allTrxs, trxs...)
	}
	return allTrxs, nil
}

// SummaryRepository methods

func (r *InMemoryRepository) SaveSummary(ctx context.Context, s *models.Summary) error {
	s.SummaryID = r.nextSummaryID
	r.nextSummaryID++
	r.summaries[s.AccountID] = s
	return nil
}

func (r *InMemoryRepository) GetSummaryByAccountID(ctx context.Context, accountID int) (*models.Summary, error) {
	summary, ok := r.summaries[accountID]
	if !ok {
		return nil, errors.New("summary not found")
	}
	return summary, nil
}

func (r *InMemoryRepository) ListSummaries(ctx context.Context) ([]*models.Summary, error) {
	var allSummaries []*models.Summary
	for _, summary := range r.summaries {
		allSummaries = append(allSummaries, summary)
	}
	return allSummaries, nil
}

// MonthSummaryRepository methods
func (r *InMemoryRepository) SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error {
	ms.MonthSummaryID = r.nextMonthSummaryID
	r.nextMonthSummaryID++
	r.monthSummaries[summaryID] = append(r.monthSummaries[ms.MonthSummaryID], ms)
	return nil
}

// GetMonthSummaryBySummaryID returns all the MonthSummary objects associated with a given Summary ID
func (r *InMemoryRepository) GetMonthSummaryBySummaryID(ctx context.Context, summaryID int) ([]*models.MonthSummary, error) {
	var monthSummaries []*models.MonthSummary
	monthSummaries = append(monthSummaries, r.monthSummaries[summaryID]...)
	return monthSummaries, nil
}
