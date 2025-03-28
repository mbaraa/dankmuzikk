package app

import (
	"dankmuzikk/app/models"
)

func (a *App) GetAccountByEmail(email string) (models.Account, error) {
	return a.repo.GetAccountByEmail(email)
}

func (a *App) GetProfileByAccountEmail(email string) (models.Profile, error) {
	account, err := a.repo.GetAccountByEmail(email)
	if err != nil {
		return models.Profile{}, err
	}

	profile, err := a.repo.GetProfileForAccount(account.Id)
	if err != nil {
		return models.Profile{}, err
	}

	profile.Account = account

	return profile, nil
}

type CreateNoOAuthUserArgs struct {
	Email string
	Name  string
}

func (a *App) CreateNoOAuthUser(args CreateNoOAuthUserArgs) (models.Profile, error) {
	profile := models.Profile{
		Account: models.Account{
			Email:   args.Email,
			IsOAuth: false,
		},
		Name:     args.Name,
		Username: args.Email,
	}

	return a.repo.CreateProfile(profile)
}

type CreateOAuthUserArgs struct {
	Email   string
	Name    string
	PfpLink string
}

func (a *App) CreateOAuthUser(args CreateOAuthUserArgs) (models.Profile, error) {
	profile := models.Profile{
		Account: models.Account{
			Email:   args.Email,
			IsOAuth: true,
		},
		Name:     args.Name,
		Username: args.Email,
		PfpLink:  args.PfpLink,
	}

	return a.repo.CreateProfile(profile)
}

func (a *App) CreateOtp(accountId uint, otp string) error {
	return a.cache.CreateOtp(accountId, otp)
}

func (a *App) GetOtpForAccount(accountId uint) (string, error) {
	otp, err := a.cache.GetOtpForAccount(accountId)
	if err != nil {
		return "", err
	}

	return otp, nil
}
