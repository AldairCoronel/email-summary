package controller

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aldaircoronel/email-summary/internal/models"
	"github.com/aldaircoronel/email-summary/internal/repository"
)

type TransactionController struct {
	repo repository.Repository
}

func NewTransactionController(repo repository.Repository) *TransactionController {
	return &TransactionController{
		repo: repo,
	}
}

func (c *TransactionController) ProcessCSVFile(ctx context.Context, filePath string) error {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Set the delimiter to comma
	reader.Comma = ','

	// Skip the first row
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to skip first row: %v", err)
	}

	// Loop through the remaining rows
	for {
		// Read the next row
		row, err := reader.Read()

		// Check for end of file
		if err == io.EOF {
			break
		}

		// Check for other errors
		if err != nil {
			return fmt.Errorf("failed to read row: %v", err)
		}

		// Parse the row values
		id, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("failed to parse ID: %v", err)
		}
		date, err := time.Parse("1/2", row[1])
		if err != nil {
			return fmt.Errorf("failed to parse date: %v", err)
		}
		amount, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %v", err)
		}
		isCredit := false
		if row[2][0] == '+' {
			isCredit = true
		}

		// Create a new transaction object
		transaction := &models.Transaction{
			Id:       id,
			Date:     date,
			Amount:   amount,
			IsCredit: isCredit,
		}

		// Save the transaction to the database
		err = c.repo.SaveTransaction(ctx, transaction)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %v", err)
		}

	}

	return nil
}

func computeSummary(transactions []*models.Transaction) (*models.Summary, error) {
	// Initialize variables to calculate summary
	var totalBalance, totalCredit, totalDebit float64
	var numOfCreditTransactions, numOfDebitTransactions int

	// Calculate summary values
	for _, trx := range transactions {
		totalBalance += trx.Amount
		if trx.IsCredit {
			totalCredit += trx.Amount
			numOfCreditTransactions++
		} else {
			totalDebit += trx.Amount
			numOfDebitTransactions++
		}
	}

	// Calculate averages
	var avgCredit, avgDebit float64
	if numOfCreditTransactions > 0 {
		avgCredit = totalCredit / float64(numOfCreditTransactions)
	}
	if numOfDebitTransactions > 0 {
		avgDebit = totalDebit / float64(numOfDebitTransactions)
	}

	// Create summary object
	summary := &models.Summary{
		TotalBalance:           totalBalance,
		NumOfCreditTansactions: numOfCreditTransactions,
		NumOfDebitTansactions:  numOfDebitTransactions,
		TotalAverageCredit:     avgCredit,
		TotalAverageDebit:      avgDebit,
	}

	return summary, nil
}

func computeMonthSummaries(transactions []*models.Transaction) ([]*models.MonthSummary, error) {
	// Create a map to group transactions by month
	transactionsByMonth := make(map[string][]*models.Transaction)

	for _, trx := range transactions {
		month := trx.Date.Format("2006-01")
		transactionsByMonth[month] = append(transactionsByMonth[month], trx)
	}

	// Compute month summaries
	monthSummaries := make([]*models.MonthSummary, 0)
	for month, monthTransactions := range transactionsByMonth {
		numOfCreditTransactions := 0
		numOfDebitTransactions := 0
		totalCredit := 0.0
		totalDebit := 0.0

		for _, trx := range monthTransactions {
			if trx.IsCredit {
				numOfCreditTransactions++
				totalCredit += trx.Amount
			} else {
				numOfDebitTransactions++
				totalDebit += trx.Amount
			}
		}

		averageCredit := 0.0
		averageDebit := 0.0

		if numOfCreditTransactions > 0 {
			averageCredit = totalCredit / float64(numOfCreditTransactions)
		}

		if numOfDebitTransactions > 0 {
			averageDebit = totalDebit / float64(numOfDebitTransactions)
		}

		monthSummary := &models.MonthSummary{
			Month:                  month,
			NumOfCreditTansactions: numOfCreditTransactions,
			NumOfDebitTansactions:  numOfDebitTransactions,
			AverageCredit:          averageCredit,
			AverageDebit:           averageDebit,
		}

		monthSummaries = append(monthSummaries, monthSummary)
	}

	return monthSummaries, nil
}

func (tc *TransactionController) GenerateEmailSummary(ctx context.Context) (*models.Summary, []*models.MonthSummary, error) {
	// Get all transactions from the repository
	transactions, err := tc.repo.ListTransactions(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list transactions: %v", err)
	}

	// Compute the summary statistics for all transactions
	summary, err := computeSummary(transactions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute summary: %v", err)
	}

	// Save the summary to the repository
	if err := tc.repo.SaveSummary(ctx, summary); err != nil {
		return nil, nil, fmt.Errorf("failed to save summary: %v", err)
	}

	// Compute the month summary statistics for each month
	monthSummaries, err := computeMonthSummaries(transactions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute month summaries: %v", err)
	}

	// Save the month summaries to the repository
	summaryID := summary.Id
	for _, monthSummary := range monthSummaries {
		if err := tc.repo.SaveMonthSummary(ctx, monthSummary, summaryID); err != nil {
			return nil, nil, fmt.Errorf("failed to save month summary: %v", err)
		}
	}

	return summary, monthSummaries, nil
}
