package main

import (
	"context"
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
	csvFilePath, _ := filepath.Abs("./sample/txns.csv")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load database connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")

	db, err := database.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set the repository to our concrete PostgreSQL repository
	repository.SetRepository(db)

	// Initialize the controller with the database as the repository
	ctrl := controller.NewTransactionController(db)

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

	// Send email summary
	to := []string{os.Getenv("EMAIL_TO")}
	subject := "Transaction Summary"
	body := view.RenderEmailBody(summary, monthSummaries)
	if err := emailService.SendEmail(to, subject, body); err != nil {
		log.Fatal(err)
	}

	// Print message when email is successfully sent
	log.Println("Email summary sent!")
}
