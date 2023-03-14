package main

import (
	"log"
	"path/filepath"

	"github.com/aldaircoronel/email-summary/internal/controller"
	"github.com/aldaircoronel/email-summary/internal/database"
	"github.com/aldaircoronel/email-summary/internal/repository"
)

func main() {
	csvFilePath, _ := filepath.Abs("./sample/txns.csv")

	// Initialize the database connection
	connStr := "postgres://aldair:hola@localhost:5432/challenge?sslmode=disable"

	db, err := database.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set the repository to our concrete PostgreSQL repository
	repository.SetRepository(db)

	// Initialize the controller with the database as the repository
	ctrl := controller.NewTransactionController(db)

	// Process the CSV file and save transactions to the database
	if err := ctrl.ProcessCSVFile(csvFilePath); err != nil {
		log.Fatal(err)
	}

	// // Generate the email summary
	// summary, err := ctrl.GenerateEmailSummary()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Send the email summary
	// if err := view.SendEmail(summary); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Email summary sent!")
}
