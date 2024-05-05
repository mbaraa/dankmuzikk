package db

import "dankmuzikk/models"

type AllowedModels interface {
	models.Account | models.Profile | models.EmailVerificationCode
	GetId() uint
}
