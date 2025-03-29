package actions

import "dankmuzikk/app/models"

type Cache interface {
	StoreLyrics(songId uint, lyrics []string) error
	GetLyrics(songId uint) ([]string, error)

	SetAuthenticatedAccount(sessionToken string, account models.Account) error
	GetAuthenticatedAccount(sessionToken string) (models.Account, error)
	InvalidateAuthenticatedAccount(sessionToken string) error
}
