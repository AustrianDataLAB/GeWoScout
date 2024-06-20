package notification

import (
	"fmt"
	"net/smtp"
	"os"
)

var (
	emailFrom = os.Getenv("NOTIFICATION_EMAIL_ADDRESS")
	emailPass = os.Getenv("NOTIFICATION_EMAIL_PASSWORD")
)

func sendHtmlEmail(to []string, subject, content string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailFrom, emailPass, "smtp.gmail.com")

	// Set up email details.
	headers := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\""
	body := content

	// Concatenate email parts.
	//msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", strings.Join(to, ", "), subject, headers, body))
	msg := []byte(fmt.Sprintf("Subject: %s\r\n%s\r\n\r\n%s", subject, headers, body))

	// Set up the SMTP server configuration.
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		emailFrom,
		to,
		msg,
	)

	// Check for errors.
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
