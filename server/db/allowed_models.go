package db

import "dankmuzikk/app/models"

type AllowedModels interface {
	models.Account | models.Profile | models.EmailVerificationCode |
		models.Song | models.Playlist | models.PlaylistSong | models.PlaylistOwner |
		models.History | models.PlaylistSongVoter
	GetId() uint
}
