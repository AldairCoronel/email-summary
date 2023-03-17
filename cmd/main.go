package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aldaircoronel/email-summary/internal/controller"
	"github.com/aldaircoronel/email-summary/internal/database"
	"github.com/aldaircoronel/email-summary/internal/repository"
	"github.com/aldaircoronel/email-summary/internal/view"
	"github.com/joho/godotenv"
)

func main() {
	// Get the file path of the input csv files.
	csvFilePath, _ := filepath.Abs("./sample/txns.csv")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load database connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")

	// Instanciate a new PostgreSQL repository
	db, err := database.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set the repository to our concrete PostgreSQL repository
	repository.SetRepository(db)

	// Initialize the controller with the database as the repository
	ctrl := controller.NewTransactionController(db)

	// Create a new account
	accountID, err := ctrl.CreateAccount(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New account created with ID: %d", accountID)

	// Store the account ID somewhere in the controller because we will need it later
	ctrl.SetAccountID(accountID)

	// Process the CSV file and save transactions to the database
	if err := ctrl.ProcessCSVFile(context.Background(), csvFilePath); err != nil {
		log.Fatal(err)
	}
	// Generate the email summary
	summary, monthSummaries, err := ctrl.GenerateEmailSummary(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Load email service configuration from environment variables
	emailCfg := &view.SMTPConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
	}

	// Initialize email service with SMTP configuration
	emailService := view.NewSMTPService(emailCfg)

	// Get EMAIL_TO flag value
	emailTo := flag.String("emailTo", "", "The email address to send the summary to")

	// Parse flags
	flag.Parse()

	to := []string{*emailTo}
	if *emailTo == "" {
		log.Fatal("The -emailTo flag is required")
	}
	subject := "Transaction Summary"
	body, err := view.RenderEmailBody(summary, monthSummaries)
	if err != nil {
		log.Fatal(err)
	}
	if err := emailService.SendEmail(to, subject, body); err != nil {
		log.Fatal(err)
	}

	// Print message when email is successfully sent
	log.Println("Email summary sent!")

}
