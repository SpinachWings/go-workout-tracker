package utils

import (
	"crypto/tls"
	"errors"
	"github.com/go-gomail/gomail"
	"os"
	"strconv"
)

func SendEmail(subject string, body string, userEmail string) error {
	from := os.Getenv("EMAIL_FROM")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return errors.New("unable to convert SMTP_PORT env var to int")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
