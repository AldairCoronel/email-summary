package view

import "net/smtp"

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
