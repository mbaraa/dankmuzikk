package actions

import "dankmuzikk/app/models"

type Cache interface {
	SetAuthenticatedAccount(sessionToken string, account models.Account) error
	GetAuthenticatedAccount(sessionToken string) (models.Account, error)
	InvalidateAuthenticatedAccount(sessionToken string) error

	SetGoogleLoginState(state string) error
	GetGoogleLoginState() (string, error)
}
