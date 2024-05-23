package db

import "dankmuzikk/models"

type AllowedModels interface {
	models.Account | models.Profile | models.EmailVerificationCode |
		models.Song | models.Playlist | models.PlaylistSong | models.PlaylistOwner |
		models.History
	GetId() uint
}
