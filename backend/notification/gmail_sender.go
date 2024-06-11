package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func GmailDemo() {
	_ = os.Setenv("GMAIL_USER", "noreply.gewoscout@gmail.com")
	//_ = os.Setenv("GMAIL_PASS", "<FILL IN HERE>")

	// Get credentials from environment variables.
	email := os.Getenv("GMAIL_USER")
	password := os.Getenv("GMAIL_PASS")

	// Set up authentication information.
	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	// Set up email details.
	from := email
	to := []string{"nikolaus.spring1@gmail.com"}
	headers := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\""
	subject := "HTML test 2"
	body :=
		`
		<!DOCTYPE html>
		<html>
		<head>
		  	<title>Your HTML Email</title>
		</head>
		<body>
		  	<h1>This is an HTML email</h1>
		  	<p>Hello, this is a test email with HTML content!</p>
		</body>
		</html>
	`

	// Concatenate email parts.
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", strings.Join(to, ", "), subject, headers, body))

	// Set up the SMTP server configuration.
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		msg,
	)

	// Check for errors.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
