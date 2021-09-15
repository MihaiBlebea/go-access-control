package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
)

//go:embed confirm.tmpl
var emailTmpl embed.FS

type Service interface {
	ConfirmEmail(email string) error
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) ConfirmEmail(email string) error {
	subject := "Please confirm your email"

	// Sender data.
	from := "mihai@gmail.com"
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Receiver email address.
	to := []string{email}

	// smtp server configuration.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Authentication.
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// t, _ := template.ParseFiles("template.html")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"ToTitle": strings.Title,
		"ToUpper": strings.ToUpper,
	}

	t, err := template.New("confirm.tmpl").Funcs(funcMap).ParseFS(emailTmpl, "confirm.tmpl")
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	if err := t.Execute(&body, struct {
		Name        string
		Message     string
		ConfirmLink string
	}{
		Name:        "Puneet Singh",
		Message:     "This is a test message in a HTML template",
		ConfirmLink: randomString(10),
	}); err != nil {
		return err
	}

	// Sending email.
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, from, to, body.Bytes()); err != nil {
		return err
	}

	return nil
}

func randomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
