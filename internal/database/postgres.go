package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aldaircoronel/email-summary/internal/models"
	_ "github.com/lib/pq"
)

// Define the PostgreSQL repository struct
type PostgresRepository struct {
	db *sql.DB
}

// Create a new PostgreSQL repository instance
func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PostgresRepository{db: db}, nil
}

// Implement the SaveAccount method of the Repository interface
func (pr *PostgresRepository) SaveAccount(ctx context.Context) (int, error) {
	// Construct the SQL query
	query := `
		INSERT INTO accounts DEFAULT VALUES
		RETURNING account_id
	`

	// Execute the query and retrieve the new account_id
	var accountID int
	err := pr.db.QueryRowContext(ctx, query).Scan(&accountID)
	if err != nil {
		return 0, err
	}

	return accountID, nil
}

// GetAccountByID retrieves the account with the given ID
func (pr *PostgresRepository) GetAccountByID(ctx context.Context, id int) (*models.Account, error) {
	query := `SELECT account_id FROM accounts WHERE account_id = $1`

	var accountID int
	err := pr.db.QueryRowContext(ctx, query, id).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("acccount with id %d not found", id)
		}
		return nil, fmt.Errorf("error getting account from database: %v", err)
	}

	return &models.Account{
		AccountID: accountID,
	}, nil
}

// Implement the SaveTransaction method of the Repository interface
func (pr *PostgresRepository) SaveTransaction(ctx context.Context, trx *models.Transaction) error {
	query := `INSERT INTO transactions (account_id, id, date, amount, is_credit) VALUES ($1, $2, $3, $4, $5)`
	_, err := pr.db.ExecContext(ctx, query, trx.AccountID, trx.ID, trx.Date, trx.Amount, trx.IsCredit)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}
	return nil
}

// Implement the GetTransactionByAccountID method of the Repository interface
func (pr *PostgresRepository) GetTransactionByAccountID(ctx context.Context, accountID int) ([]*models.Transaction, error) {
	query := `SELECT transaction_id, account_id, id, date, amount, is_credit FROM transactions WHERE account_id=$1`
	rows, err := pr.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.TransactionID, &transaction.AccountID, &transaction.ID, &transaction.Date, &transaction.Amount, &transaction.IsCredit); err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read transaction rows: %v", err)
	}
	return transactions, nil
}

// Implement the ListTransactions method of the Repository interface
func (pr *PostgresRepository) ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	query := `SELECT transaction_id, account_id, id, date, amount, is_credit FROM transactions ORDER BY date DESC`
	rows, err := pr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		trx := &models.Transaction{}
		if err := rows.Scan(&trx.TransactionID, &trx.AccountID, &trx.ID, &trx.Date, &trx.Amount, &trx.IsCredit); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, trx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate transactions: %v", err)
	}

	return transactions, nil
}

// Implement the SaveSummary method of the Repository interface
func (pr *PostgresRepository) SaveSummary(ctx context.Context, s *models.Summary) error {
	query := `
		INSERT INTO summary (
			account_id,
			total_balance, 
			total_transactions, 
			num_of_credit_transactions, 
			num_of_debit_transactions, 
			total_average_credit, 
			total_average_debit
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING summary_id
	`
	row := pr.db.QueryRowContext(
		ctx,
		query,
		s.AccountID,
		s.TotalBalance,
		s.TotalTransactions,
		s.NumOfCreditTransactions,
		s.NumOfDebitTransactions,
		s.TotalAverageCredit,
		s.TotalAverageDebit,
	)
	if err := row.Scan(&s.SummaryID); err != nil {
		return fmt.Errorf("failed to save summary: %v", err)
	}
	return nil
}

// Implement the GetSummaryByAccountID method of the Repository interface
func (pr *PostgresRepository) GetSummaryByAccountID(ctx context.Context, accountID int) (*models.Summary, error) {
	query := `
		SELECT summary_id, total_balance, total_transactions, num_of_credit_transactions, num_of_debit_transactions, total_average_credit, total_average_debit
		FROM summary
		WHERE account_id = $1
	`
	row := pr.db.QueryRowContext(ctx, query, accountID)

	summary := &models.Summary{}
	err := row.Scan(&summary.SummaryID, &summary.TotalBalance, &summary.TotalTransactions, &summary.NumOfCreditTransactions, &summary.NumOfDebitTransactions, &summary.TotalAverageCredit, &summary.TotalAverageDebit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get summary by id: %v", err)
		}
		return nil, fmt.Errorf("failed to get summary by account ID: %v", err)
	}

	return summary, nil
}

// ListSummaries returns a list of all summaries for all accounts.
func (pr *PostgresRepository) ListSummaries(ctx context.Context) ([]*models.Summary, error) {
	query := `
		SELECT summary_id, account_id, total_balance, total_transactions, num_of_credit_transactions, num_of_debit_transactions, total_average_credit, total_average_debit 
		FROM summary
	`
	rows, err := pr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list summaries: %v", err)
	}
	defer rows.Close()

	summaries := make([]*models.Summary, 0)
	for rows.Next() {
		var summary models.Summary
		err = rows.Scan(
			&summary.SummaryID,
			&summary.AccountID,
			&summary.TotalBalance,
			&summary.TotalTransactions,
			&summary.NumOfCreditTransactions,
			&summary.NumOfDebitTransactions,
			&summary.TotalAverageCredit,
			&summary.TotalAverageDebit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan summary row: %v", err)
		}
		summaries = append(summaries, &summary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list summaries: %v", err)
	}

	return summaries, nil
}

// Implement the SaveMonthSummary method of the Repository interface
func (pr *PostgresRepository) SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error {
	query := `
		INSERT INTO month_summary (
			month, 
			total_balance, 
			total_transactions, 
			num_of_credit_transactions, 
			num_of_debit_transactions, 
			average_credit, 
			average_debit, 
			summary_id
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := pr.db.ExecContext(
		ctx,
		query,
		ms.Month,
		ms.TotalBalance,
		ms.TotalTransactions,
		ms.NumOfCreditTransactions,
		ms.NumOfDebitTransactions,
		ms.AverageCredit,
		ms.AverageDebit,
		summaryID,
	)
	if err != nil {
		return fmt.Errorf("failed to save month summary: %v", err)
	}
	return nil
}

// GetMonthSummaryBySummaryID returns a month summary by summary id
func (pr *PostgresRepository) GetMonthSummaryBySummaryID(ctx context.Context, summaryID int) ([]*models.MonthSummary, error) {
	query := `
		SELECT 
			month_summary_id, 
			month, 
			total_balance, 
			total_transactions, 
			num_of_credit_transactions, 
			num_of_debit_transactions, 
			average_credit, 
			average_debit,
			summary_id
		FROM month_summary
		WHERE summary_id = $1
	`
	rows, err := pr.db.QueryContext(ctx, query, summaryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get month summary by summary id: %v", err)
	}
	defer rows.Close()

	var monthSummaries []*models.MonthSummary
	for rows.Next() {
		ms := new(models.MonthSummary)
		err := rows.Scan(
			&ms.MonthSummaryID,
			&ms.Month,
			&ms.TotalBalance,
			&ms.TotalTransactions,
			&ms.NumOfCreditTransactions,
			&ms.NumOfDebitTransactions,
			&ms.AverageCredit,
			&ms.AverageDebit,
			&ms.SummaryID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan month summary: %v", err)
		}
		monthSummaries = append(monthSummaries, ms)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get month summary by summary id: %v", err)
	}

	return monthSummaries, nil
}

// This function closes the database connection by calling the Close() function on the database object.
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
