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

func (c *TransactionController) ProcessCSVFile(filePath string) error {
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
		err = c.repo.SaveTransaction(context.Background(), transaction)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %v", err)
		}

	}

	return nil
}
