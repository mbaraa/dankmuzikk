package actions

import "dankmuzikk/app/models"

type Mailer interface {
	SendOtpEmail(profile models.Profile, code string) error
}
