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

// TransactionController defines a controller for handling transactions.
type TransactionController struct {
	repo      repository.Repository
	accountID int
}

// NewTransactionController creates a new instance of TransactionController.
func NewTransactionController(repo repository.Repository) *TransactionController {
	return &TransactionController{
		repo: repo,
	}
}

// CreateAccount creates a new account and returns its ID
func (c *TransactionController) CreateAccount(ctx context.Context) (int, error) {
	accountID, err := c.repo.SaveAccount(ctx)
	if err != nil {
		return 0, fmt.Errorf("error creating account: %v", err)
	}
	return accountID, nil
}

// SetAccountID stores the account ID in the controller
func (c *TransactionController) SetAccountID(accountID int) {
	c.accountID = accountID
}

// ProcessCSVFile reads a CSV file from the given file path and inserts its contents into the database
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
			ID:        id,
			Date:      date,
			Amount:    amount,
			IsCredit:  isCredit,
			AccountID: c.accountID,
		}

		// Save the transaction to the database
		err = c.repo.SaveTransaction(ctx, transaction)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %v", err)
		}
	}

	return nil
}

/*
This function takes a slice of *models.Transaction and returns a pointer to models.Summary and an error. It computes summary statistics for all transactions in the slice, including total balance, total transactions, number of credit and debit transactions, and average credit and debit amounts. If successful, it returns a pointer to the computed models.Summary and a nil error. If there was an error, it returns a nil pointer and an error.
*/
func ComputeSummary(transactions []*models.Transaction) (*models.Summary, error) {
	var totalBalance, totalCredit, totalDebit float64
	var totalTransactions, numCreditTransactions, numDebitTransactions int

	for _, transaction := range transactions {
		if transaction.IsCredit {
			totalBalance += transaction.Amount
			totalCredit += transaction.Amount
			numCreditTransactions++
		} else {
			totalBalance -= transaction.Amount
			totalDebit += transaction.Amount
			numDebitTransactions++
		}
		totalTransactions++
	}

	totalAverageCredit := 0.0
	if numCreditTransactions > 0 {
		totalAverageCredit = totalCredit / float64(numCreditTransactions)
	}

	totalAverageDebit := 0.0
	if numDebitTransactions > 0 {
		totalAverageDebit = totalDebit / float64(numDebitTransactions)
	}

	summary := &models.Summary{
		AccountID:               transactions[0].AccountID,
		TotalBalance:            totalBalance,
		TotalTransactions:       totalTransactions,
		NumOfCreditTransactions: numCreditTransactions,
		NumOfDebitTransactions:  numDebitTransactions,
		TotalAverageCredit:      totalAverageCredit,
		TotalAverageDebit:       totalAverageDebit,
	}

	return summary, nil
}

/*
This function takes a slice of *models.Transaction and returns a slice of *models.MonthSummary and an error. It computes summary statistics for each month in the transactions, including total balance, total transactions, number of credit and debit transactions, and average credit and debit amounts. If successful, it returns a slice of pointers to the computed models.MonthSummary structs and a nil error. If there was an error, it returns a nil slice and an error.
*/
func ComputeMonthSummaries(transactions []*models.Transaction) ([]*models.MonthSummary, error) {
	// Create a map to hold the month summaries
	monthSummaries := make(map[string]*models.MonthSummary)

	// Iterate over the transactions and add them to the month summaries
	for _, transaction := range transactions {
		month := transaction.Date.Month().String() // Get the month name

		// Check if a month summary already exists for this month
		if _, ok := monthSummaries[month]; !ok {
			// Create a new month summary if one doesn't exist
			monthSummaries[month] = &models.MonthSummary{
				Month: month,
			}
		}

		// Add the transaction to the appropriate month summary
		monthSummary := monthSummaries[month]
		monthSummary.TotalTransactions++
		if transaction.IsCredit {
			monthSummary.TotalBalance += transaction.Amount
			monthSummary.NumOfCreditTransactions++
			monthSummary.AverageCredit += transaction.Amount
		} else {
			monthSummary.TotalBalance -= transaction.Amount
			monthSummary.NumOfDebitTransactions++
			monthSummary.AverageDebit += transaction.Amount
		}
	}

	// Calculate the averages for each month summary
	for _, monthSummary := range monthSummaries {
		if monthSummary.NumOfCreditTransactions > 0 {
			monthSummary.AverageCredit /= float64(monthSummary.NumOfCreditTransactions)
		}
		if monthSummary.NumOfDebitTransactions > 0 {
			monthSummary.AverageDebit /= float64(monthSummary.NumOfDebitTransactions)
		}
	}

	// Convert the map to a slice of month summaries and return it
	result := make([]*models.MonthSummary, 0, len(monthSummaries))
	for _, monthSummary := range monthSummaries {
		result = append(result, monthSummary)
	}
	return result, nil
}

/*
GenerateEmailSummary is a method of TransactionController that takes a context and returns a pointer to models.Summary, a slice of pointers to models.MonthSummary, and an error. It first retrieves all transactions for the account ID associated with the TransactionController instance from the repository, then computes summary statistics and month-wise summary statistics using helper functions computeSummary and computeMonthSummaries, respectively. It saves the computed summary and month summaries to the repository and returns them along with a nil error. If there was an error retrieving or computing the summary or saving the summary to the repository, it returns nil pointers and an error.
*/
func (tc *TransactionController) GenerateEmailSummary(ctx context.Context) (*models.Summary, []*models.MonthSummary, error) {
	// Get the account ID from the controller
	accountID := tc.accountID

	// Get all transactions for the account from the repository
	transactions, err := tc.repo.GetTransactionByAccountID(ctx, accountID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transactions for account %d: %v", accountID, err)
	}

	// Compute the summary statistics for all transactions
	summary, err := ComputeSummary(transactions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute summary: %v", err)
	}

	// Save the summary to the repository
	if err := tc.repo.SaveSummary(ctx, summary); err != nil {
		return nil, nil, fmt.Errorf("failed to save summary: %v", err)
	}

	// Compute the month summary statistics for each month
	monthSummaries, err := ComputeMonthSummaries(transactions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute month summaries: %v", err)
	}

	// Save the month summaries to the repository
	summaryID := summary.SummaryID
	for _, monthSummary := range monthSummaries {
		if err := tc.repo.SaveMonthSummary(ctx, monthSummary, summaryID); err != nil {
			return nil, nil, fmt.Errorf("failed to save month summary: %v", err)
		}
	}

	return summary, monthSummaries, nil
}
