package actions

import "dankmuzikk/app/models"

type Cache interface {
	StoreLyrics(songId uint, lyrics []string) error
	GetLyrics(songId uint) ([]string, error)

	SetAuthenticatedUser(sessionToken string, profile models.Profile) error
	GetAuthenticatedUser(sessionToken string) (models.Profile, error)
	InvalidateAuthenticatedUser(sessionToken string) error
}
