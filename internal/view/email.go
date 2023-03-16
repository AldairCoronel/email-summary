package view

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/aldaircoronel/email-summary/internal/models"
)

// Interface that represents an email service.
type EmailService interface {
	SendEmail(to []string, subject string, body string) error
}

// SMTPConfig contains configuration options for the SMTP service.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// SMTPService is the implementation of the EmailService interface that sends email through SMTP
type SMTPService struct {
	smtpHost string
	smtpPort string
	auth     smtp.Auth
	from     string
}

// Constructor that creates a new SMTPService
func NewSMTPService(cfg *SMTPConfig) *SMTPService {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &SMTPService{
		smtpHost: cfg.Host,
		smtpPort: cfg.Port,
		auth:     auth,
		from:     cfg.From,
	}
}

func RenderEmailBody(summary *models.Summary, monthSummaries []*models.MonthSummary) (string, error) {
	tmpl, err := template.ParseFiles("internal/view/email-template.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %v", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, struct {
		Summary        *models.Summary
		MonthSummaries []*models.MonthSummary
	}{
		Summary:        summary,
		MonthSummaries: monthSummaries,
	}); err != nil {
		return "", fmt.Errorf("failed to execute email template: %v", err)
	}

	return body.String(), nil
}

// SendEmail sends an email through SMTP
func (s *SMTPService) SendEmail(to []string, subject string, body string) error {
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body + "\r\n")

	addr := s.smtpHost + ":" + s.smtpPort
	if err := smtp.SendMail(addr, s.auth, s.from, to, msg); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
