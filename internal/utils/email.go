package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPasswordResetEmail(email, token string) error {
	// Load SMTP settings from environment variables
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	// Email content
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password:\n\nhttps://example.com/reset-password?token=%s", token)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// SMTP authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	// Sending email
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUser, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("Password reset email sent to %s\n", email)
	return nil
}
