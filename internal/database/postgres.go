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

// Implement the SaveTransaction method of the Repository interface
func (pr *PostgresRepository) SaveTransaction(ctx context.Context, trx *models.Transaction) error {
	query := `INSERT INTO transactions (id, date, amount, is_credit) VALUES ($1, $2, $3, $4)`
	_, err := pr.db.ExecContext(ctx, query, trx.Id, trx.Date, trx.Amount, trx.IsCredit)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}
	return nil
}

// Implement the GetTransactionByID method of the Repository interface
func (pr *PostgresRepository) GetTransactionByID(ctx context.Context, id int) (*models.Transaction, error) {
	query := `SELECT id, date, amount, is_credit FROM transactions WHERE id=$1`
	row := pr.db.QueryRowContext(ctx, query, id)

	trx := &models.Transaction{}
	err := row.Scan(&trx.Id, &trx.Date, &trx.Amount, &trx.IsCredit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	return trx, nil
}

// Implement the ListTransactions method of the Repository interface
func (pr *PostgresRepository) ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, date, amount, is_credit FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		trx := new(models.Transaction)
		err := rows.Scan(&trx.Id, &trx.Date, &trx.Amount, &trx.IsCredit)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, trx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

// Implement the SaveSummary method of the Repository interface
func (pr *PostgresRepository) SaveSummary(ctx context.Context, s *models.Summary) error {
	query := `INSERT INTO summary (total_balance, num_of_credit_tansactions, num_of_debit_tansactions, total_average_credit, total_average_debit)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := pr.db.QueryRowContext(ctx, query, s.TotalBalance, s.NumOfCreditTansactions, s.NumOfDebitTansactions, s.TotalAverageCredit, s.TotalAverageDebit).Scan(&s.Id)
	if err != nil {
		return err
	}
	return nil
}

// Implement the GetSummaryByID method of the Repository interface
func (pr *PostgresRepository) GetSummaryByID(ctx context.Context, id int) (*models.Summary, error) {
	s := new(models.Summary)
	query := "SELECT id, total_balance, num_of_credit_tansactions, num_of_debit_tansactions, total_average_credit, total_average_debit FROM summary WHERE id=$1"
	err := pr.db.QueryRowContext(ctx, query, id).Scan(&s.Id, &s.TotalBalance, &s.NumOfCreditTansactions, &s.NumOfDebitTansactions, &s.TotalAverageCredit, &s.TotalAverageDebit)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Implement the SaveMonthSummary method of the Repository interface
func (pr *PostgresRepository) SaveMonthSummary(ctx context.Context, ms *models.MonthSummary, summaryID int) error {
	sqlStatement := `
		INSERT INTO month_summary (month, num_of_credit_tansactions, num_of_debit_tansactions, average_credit, average_debit, summary_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	id := 0
	err := pr.db.QueryRowContext(ctx, sqlStatement,
		ms.Month,
		ms.NumOfCreditTansactions,
		ms.NumOfDebitTansactions,
		ms.AverageCredit,
		ms.AverageDebit,
		summaryID,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to save month summary: %v", err)
	}

	ms.Id = id
	return nil
}

// Implement the GetMonthSummaryByID method of the Repository interface
func (pr *PostgresRepository) GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error) {
	sqlStatement := `SELECT id, month, num_of_credit_tansactions, num_of_debit_tansactions, average_credit, average_debit, summary_id FROM month_summary WHERE id=$1`
	row := pr.db.QueryRowContext(ctx, sqlStatement, id)

	ms := &models.MonthSummary{}
	err := row.Scan(
		&ms.Id,
		&ms.Month,
		&ms.NumOfCreditTansactions,
		&ms.NumOfDebitTansactions,
		&ms.AverageCredit,
		&ms.AverageDebit,
		&ms.SummaryId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("month summary not found")
		}
		return nil, fmt.Errorf("failed to get month summary: %v", err)
	}

	return ms, nil
}

// Implement the GetMonthSummaryBySummaryID method of the Repository interface
func (pr *PostgresRepository) GetMonthSummaryBySummaryID(ctx context.Context, summaryId int) ([]*models.MonthSummary, error) {
	sqlStatement := `SELECT id, month, num_of_credit_tansactions, num_of_debit_tansactions, average_credit, average_debit, summary_id FROM month_summary WHERE summary_id=$1`
	rows, err := pr.db.QueryContext(ctx, sqlStatement, summaryId)
	if err != nil {
		return nil, fmt.Errorf("failed to get month summaries: %v", err)
	}
	defer rows.Close()

	var monthSummaries []*models.MonthSummary
	for rows.Next() {
		ms := &models.MonthSummary{}
		err := rows.Scan(
			&ms.Id,
			&ms.Month,
			&ms.NumOfCreditTansactions,
			&ms.NumOfDebitTansactions,
			&ms.AverageCredit,
			&ms.AverageDebit,
			&ms.SummaryId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get month summaries: %v", err)
		}

		monthSummaries = append(monthSummaries, ms)
	}

	return monthSummaries, nil
}

// Implement the ListMonthSummaries method of the Repository interface
func (pr *PostgresRepository) ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, month, num_of_credit_transactions, num_of_debit_transactions, average_credit, average_debit, summary_id FROM month_summaries")
	if err != nil {
		return nil, fmt.Errorf("failed to list month summaries: %v", err)
	}
	defer rows.Close()

	var summaries []*models.MonthSummary
	for rows.Next() {
		var summary models.MonthSummary
		err = rows.Scan(&summary.Id, &summary.Month, &summary.NumOfCreditTansactions, &summary.NumOfDebitTansactions, &summary.AverageCredit, &summary.AverageDebit, &summary.SummaryId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan month summary: %v", err)
		}
		summaries = append(summaries, &summary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list month summaries: %v", err)
	}

	return summaries, nil
}

// This function closes the database connection by calling the Close() function on the database object.
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
