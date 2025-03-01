package mailer

import (
	"bytes"
	"context"
	"dankmuzikk/views/emails"
)

func SendOtpEmail(name, email, code string) error {
	buf := bytes.NewBuffer([]byte{})
	err := emails.OtpEmail(name, code).Render(context.Background(), buf)
	if err != nil {
		return err
	}
	return sendEmail("Email verification", buf.String(), email)
}
