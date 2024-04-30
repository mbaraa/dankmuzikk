package db

import "dankmuzikk/models"

type AllowedModels interface {
	models.Account | models.Profile
	GetId() uint
}
