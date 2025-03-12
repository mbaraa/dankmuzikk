package mailer

import (
	"bytes"
	"context"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"fmt"
	"net/smtp"
)

type SmtpMailer struct {
}

func New() *SmtpMailer {
	return &SmtpMailer{}
}

func (m *SmtpMailer) SendOtpEmail(profile models.Profile, code string) error {
	buf := bytes.NewBuffer([]byte{})
	err := otpEmail(profile.Name, code).Render(context.Background(), buf)
	if err != nil {
		return err
	}

	return sendEmail("Email verification", buf.String(), profile.Account.Email)
}

func sendEmail(subject, content, to string) error {
	receiver := []string{to}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_subject := "Subject: " + subject
	_to := "To: " + to
	_from := fmt.Sprintf("From: Baraa from DankMuzikk <%s>", config.Env().Smtp.Username)
	body := []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s", _from, _to, _subject, mime, content))

	addr := config.Env().Smtp.Host + ":" + config.Env().Smtp.Port
	auth := smtp.PlainAuth("", config.Env().Smtp.Username, config.Env().Smtp.Password, config.Env().Smtp.Host)

	return smtp.SendMail(addr, auth, config.Env().Smtp.Username, receiver, body)
}
