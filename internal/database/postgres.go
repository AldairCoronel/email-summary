package database

import (
	"context"
	"database/sql"
	"encoding/json"
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
	sql := `
		INSERT INTO transactions (id, date, amount, is_credit)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			date = excluded.date,
			amount = excluded.amount,
			is_credit = excluded.is_credit
	`

	_, err := pr.db.ExecContext(ctx, sql, trx.Id, trx.Date, trx.Amount, trx.IsCredit)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	return nil
}

// Implement the GetTransactionByID method of the Repository interface
func (pr *PostgresRepository) GetTransactionByID(ctx context.Context, id int) (*models.Transaction, error) {
	sql := `
		SELECT id, date, amount, is_credit
		FROM transactions
		WHERE id = $1
	`

	row := pr.db.QueryRowContext(ctx, sql, id)
	trx := &models.Transaction{}

	err := row.Scan(&trx.Id, &trx.Date, &trx.Amount, &trx.IsCredit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID: %v", err)
	}

	return trx, nil
}

// Implement the ListTransactions method of the Repository interface
func (pr *PostgresRepository) ListTransactions(ctx context.Context) ([]*models.Transaction, error) {
	sql := `
		SELECT id, date, amount, is_credit
		FROM transactions
		ORDER BY date DESC
	`

	rows, err := pr.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %v", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		trx := &models.Transaction{}
		err = rows.Scan(&trx.Id, &trx.Date, &trx.Amount, &trx.IsCredit)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, trx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list transactions: %v", err)
	}

	return transactions, nil
}

// Implement the SaveMonthSummary method of the Repository interface
func (pr *PostgresRepository) SaveMonthSummary(ctx context.Context, ms *models.MonthSummary) error {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	// Check if the summary exists
	summaryID := 0
	row := tx.QueryRowContext(ctx, "SELECT id FROM summary ORDER BY id DESC LIMIT 1")
	err = row.Scan(&summaryID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting summary ID: %v", err)
	}

	// Save month summary
	_, err = tx.ExecContext(ctx, "INSERT INTO month_summary (month, year, num_of_transactions, average_credit, average_debit, summary_id) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (month) DO UPDATE SET year=EXCLUDED.year, num_of_transactions=EXCLUDED.num_of_transactions, average_credit=EXCLUDED.average_credit, average_debit=EXCLUDED.average_debit, summary_id=EXCLUDED.summary_id",
		ms.Month, ms.Year, ms.NumOfTransactions, ms.AverageCredit, ms.AverageDebit, summaryID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error saving month summary: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

// Implement the GetMonthSummaryByID method of the Repository interface
func (p *PostgresRepository) GetMonthSummaryByID(ctx context.Context, id int) (*models.MonthSummary, error) {
	monthSummary := &models.MonthSummary{}
	query := `SELECT month, year, num_of_transactions, average_credit, average_debit, summary_id 
	          FROM month_summary WHERE month = $1`
	row := p.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&monthSummary.Month, &monthSummary.Year, &monthSummary.NumOfTransactions,
		&monthSummary.AverageCredit, &monthSummary.AverageDebit, &monthSummary.SummaryId)
	if err != nil {
		return nil, err
	}
	return monthSummary, nil
}

// Implement the ListMonthSummaries method of the Repository interface
func (r *PostgresRepository) ListMonthSummaries(ctx context.Context) ([]*models.MonthSummary, error) {
	query := `SELECT month, year, num_of_transactions, average_credit, average_debit, summary_id
              FROM month_summary`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	monthSummaries := make([]*models.MonthSummary, 0)

	for rows.Next() {
		ms := new(models.MonthSummary)

		var year, numOfTransactions int
		var avgCredit, avgDebit float64
		var summaryID int64

		err = rows.Scan(&ms.Month, &year, &numOfTransactions, &avgCredit, &avgDebit, &summaryID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		ms.Year = int(year)
		ms.NumOfTransactions = int(numOfTransactions)
		ms.AverageCredit = avgCredit
		ms.AverageDebit = avgDebit
		ms.SummaryId = int(summaryID)

		monthSummaries = append(monthSummaries, ms)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return monthSummaries, nil
}

func (p *PostgresRepository) SaveSummary(ctx context.Context, s *models.Summary) error {
	// Convert monthly summaries to JSONB
	monthlySummaries, err := json.Marshal(s.MonthlySummaries)
	if err != nil {
		return fmt.Errorf("failed to marshal monthly summaries to JSONB: %w", err)
	}

	// Begin transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert summary row
	res, err := tx.ExecContext(ctx, "INSERT INTO summary (total_balance, num_of_total_transactions, monthly_summaries) VALUES ($1, $2, $3)",
		s.TotalBalance, s.NumOfTotalTransactions, monthlySummaries)
	if err != nil {
		return fmt.Errorf("failed to insert summary row: %w", err)
	}

	// Get summary ID
	summaryID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get summary ID: %w", err)
	}

	// Insert month summary rows
	for _, ms := range s.MonthlySummaries {
		_, err = tx.ExecContext(ctx, "INSERT INTO month_summary (month, year, num_of_transactions, average_credit, average_debit, summary_id) VALUES ($1, $2, $3, $4, $5, $6)",
			ms.Month, ms.Year, ms.NumOfTransactions, ms.AverageCredit, ms.AverageDebit, summaryID)
		if err != nil {
			return fmt.Errorf("failed to insert month summary row: %w", err)
		}
	}

	return nil
}

// Implement the GetSummaryByID method of the Repository interface
func (r *PostgresRepository) GetSummaryByID(ctx context.Context, id int) (*models.Summary, error) {
	// Construct the SQL query
	query := `SELECT id, total_balance, num_of_total_transactions, monthly_summaries FROM summary WHERE id = $1`

	// Execute the query and retrieve the row
	row := r.db.QueryRowContext(ctx, query, id)

	// Create a new Summary instance to hold the result
	summary := &models.Summary{}

	// Scan the row and populate the Summary instance
	if err := row.Scan(&summary.Id, &summary.TotalBalance, &summary.NumOfTotalTransactions, &summary.MonthlySummaries); err != nil {
		// Return the error
		return nil, fmt.Errorf("error getting summary by ID: %v", err)
	}

	return summary, nil
}

// This function closes the database connection by calling the Close() function on the database object.
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
