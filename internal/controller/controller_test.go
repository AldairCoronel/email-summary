package controller_test

import (
	"context"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aldaircoronel/email-summary/internal/controller"
	"github.com/aldaircoronel/email-summary/internal/models"
	"github.com/aldaircoronel/email-summary/internal/repository"
)

func TestTransactionController_CreateAccount(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	ctrl := controller.NewTransactionController(repo)

	accountID, err := ctrl.CreateAccount(context.Background())
	if err != nil {
		t.Errorf("expected error to be nil, but got %v", err)
	}

	if accountID <= 0 {
		t.Errorf("expected account ID to be greater than 0, but got %d", accountID)
	}
}

func TestTransactionController_ProcessCSVFile(t *testing.T) {
	ctx := context.Background()
	// Create an in-memory repository
	repo := repository.NewInMemoryRepository()

	// Create a new transaction controller with the in-memory repository
	controller := controller.NewTransactionController(repo)

	// Set AccountID
	controller.SetAccountID(1)

	// Create a temporary CSV file with test data
	testData := []string{
		"ID,Date,Amount",
		"1,1/1,100.00",
		"2,1/2,50.00",
		"3,1/3,-75.00",
		"4,1/4,+125.00",
	}
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(strings.Join(testData, "\n"))); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call the ProcessCSVFile method with the temporary file
	err = controller.ProcessCSVFile(ctx, tmpfile.Name())
	if err != nil {
		t.Fatalf("ProcessCSVFile failed with error: %v", err)
	}

	// Verify that the transactions were saved to the repository
	transactions, err := repo.GetTransactionByAccountID(ctx, 1)
	if err != nil {
		t.Fatalf("GetTransactionsByAccountID failed with error: %v", err)
	}

	if len(transactions) != 4 {
		t.Errorf("Expected 4 transactions, but got %d", len(transactions))
	}

	expectedIDs := []int{1, 2, 3, 4}
	for i, transaction := range transactions {
		if transaction.ID != expectedIDs[i] {
			t.Errorf("Expected transaction with ID %d, but got %d", expectedIDs[i], transaction.ID)
		}
	}
}

func TestComputeSummary(t *testing.T) {
	// Create a slice of transactions with known values
	transactions := []*models.Transaction{
		{ID: 1, AccountID: 1, IsCredit: true, Amount: 100},
		{ID: 2, AccountID: 1, IsCredit: true, Amount: 200},
		{ID: 3, AccountID: 1, IsCredit: false, Amount: -50},
		{ID: 4, AccountID: 1, IsCredit: false, Amount: -75},
	}

	// Compute the summary
	summary, err := controller.ComputeSummary(transactions)

	// Check if there was an error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Create the expected summary
	expectedSummary := &models.Summary{
		AccountID:               1,
		TotalBalance:            425,
		TotalTransactions:       4,
		NumOfCreditTransactions: 2,
		NumOfDebitTransactions:  2,
		TotalAverageCredit:      150,
		TotalAverageDebit:       -62.5,
	}

	// Compare the computed summary with the expected summary
	if !reflect.DeepEqual(summary, expectedSummary) {
		t.Errorf("unexpected summary value: got %+v, expected %+v", summary, expectedSummary)
	}
}

func sortByMonth(monthSummaries []*models.MonthSummary) {
	sort.Slice(monthSummaries, func(i, j int) bool {
		m1, _ := time.Parse("January", monthSummaries[i].Month)
		m2, _ := time.Parse("January", monthSummaries[j].Month)
		return m1.Before(m2)
	})
}

func TestComputeMonthSummaries(t *testing.T) {
	transactions := []*models.Transaction{
		{ID: 1, AccountID: 123, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), Amount: 100, IsCredit: true},
		{ID: 2, AccountID: 123, Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC), Amount: -50, IsCredit: false},
		{ID: 3, AccountID: 123, Date: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC), Amount: 75, IsCredit: true},
		{ID: 4, AccountID: 123, Date: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC), Amount: -25, IsCredit: false},
		{ID: 5, AccountID: 123, Date: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC), Amount: 200.0, IsCredit: true},
	}

	expectedSummaries := []*models.MonthSummary{
		{Month: "January", TotalTransactions: 2, TotalBalance: 150.0, NumOfCreditTransactions: 1, NumOfDebitTransactions: 1, AverageCredit: 100.0, AverageDebit: -50.0},
		{Month: "February", TotalTransactions: 2, TotalBalance: 100.0, NumOfCreditTransactions: 1, NumOfDebitTransactions: 1, AverageCredit: 75.0, AverageDebit: -25.0},
		{Month: "March", TotalTransactions: 1, TotalBalance: 200.0, NumOfCreditTransactions: 1, NumOfDebitTransactions: 0, AverageCredit: 200.0, AverageDebit: 0.0},
	}

	// Call the function to get the actual result
	actualSummaries, _ := controller.ComputeMonthSummaries(transactions)
	sortByMonth(actualSummaries)

	// Check the length of the actual result
	if len(actualSummaries) != len(expectedSummaries) {
		t.Fatalf("Unexpected number of month summaries. Expected %d, but got %d", len(expectedSummaries), len(actualSummaries))
	}

	// Check the values of each month summary
	for i, expected := range expectedSummaries {
		actual := actualSummaries[i]
		if expected.Month != actual.Month {
			t.Errorf("Unexpected month. Expected %s, but got %s", expected.Month, actual.Month)
		}
		if expected.TotalTransactions != actual.TotalTransactions {
			t.Errorf("Unexpected total transactions. Expected %d, but got %d", expected.TotalTransactions, actual.TotalTransactions)
		}
		if expected.TotalBalance != actual.TotalBalance {
			t.Errorf("Unexpected total balance. Expected %f, but got %f", expected.TotalBalance, actual.TotalBalance)
		}
		if expected.NumOfCreditTransactions != actual.NumOfCreditTransactions {
			t.Errorf("Unexpected number of credit transactions. Expected %d, but got %d", expected.NumOfCreditTransactions, actual.NumOfCreditTransactions)
		}
		if expected.NumOfDebitTransactions != actual.NumOfDebitTransactions {
			t.Errorf("Unexpected number of debit transactions. Expected %d, but got %d", expected.NumOfDebitTransactions, actual.NumOfDebitTransactions)
		}
		if expected.AverageCredit != actual.AverageCredit {
			t.Errorf("Unexpected average credit. Expected %f, but got %f", expected.AverageCredit, actual.AverageCredit)
		}
		if expected.AverageDebit != actual.AverageDebit {
			t.Errorf("Unexpected average debit. Expected %f, but got %f", expected.AverageDebit, actual.AverageDebit)
		}
	}
}

func TestTransactionController_GenerateEmailSummary(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository
	repo := repository.NewInMemoryRepository()

	// Create a new transaction controller with a mock account ID and the in-memory repository
	accountID := 1
	// Create a new transaction controller with the in-memory repository
	controller := controller.NewTransactionController(repo)

	// Set AccountID
	controller.SetAccountID(accountID)

	// Add some transactions to the repository
	transactions := []*models.Transaction{
		{
			ID:        1,
			AccountID: accountID,
			IsCredit:  false,
			Amount:    -100.0,
			Date:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			AccountID: accountID,
			IsCredit:  true,
			Amount:    50.0,
			Date:      time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        3,
			AccountID: accountID,
			IsCredit:  true,
			Amount:    75.0,
			Date:      time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, transaction := range transactions {
		if err := repo.SaveTransaction(ctx, transaction); err != nil {
			t.Fatalf("failed to save transaction: %v", err)
		}
	}

	// Generate the email summary using the transaction controller
	summary, _, err := controller.GenerateEmailSummary(ctx)
	if err != nil {
		t.Fatalf("failed to generate email summary: %v", err)
	}

	// Check that the summary was saved to the repository
	savedSummary, err := repo.GetSummaryByAccountID(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get saved summary: %v", err)
	}
	if !reflect.DeepEqual(summary, savedSummary) {
		t.Errorf("saved summary does not match expected. Expected %+v, but got %+v", summary, savedSummary)
	}

}
