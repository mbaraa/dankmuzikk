package app

import (
	"dankmuzikk/app/models"
	"time"
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

func (a *App) CreateOtp(otp models.EmailVerificationCode) error {
	// TODO: use cache instead.
	return a.repo.CreateOtp(otp)
}

func (a *App) GetOtpForAccount(accountId uint) (models.EmailVerificationCode, error) {
	otp, err := a.repo.GetOtpForAccount(accountId)
	if err != nil {
		return models.EmailVerificationCode{}, err
	}

	if otp.CreatedAt.Add(time.Hour / 2).Before(time.Now()) {
		return models.EmailVerificationCode{}, &ErrExpiredVerificationCode{}
	}

	_ = a.repo.DeleteOtpsForAccount(accountId)

	return otp, nil
}
