package login

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"time"
)

type EmailLoginService struct {
	accountRepo db.CRUDRepo[models.Account]
	profileRepo db.CreatorRepo[models.Profile]
	jwtUtil     jwt.Manager[any]
}

func NewEmailLoginService(
	accountRepo db.CRUDRepo[models.Account],
	profileRepo db.CreatorRepo[models.Profile],
	jwtUtil jwt.Manager[any],
) *EmailLoginService {
	return &EmailLoginService{accountRepo, profileRepo, jwtUtil}
}

func (e *EmailLoginService) Login(user entities.LoginRequest) (string, error) {
	return "", nil
}

func (e *EmailLoginService) Signup(user entities.SignupRequest) (string, error) {
	// creating a new account will create the account underneeth it.
	err := e.profileRepo.Add(&models.Profile{
		Account: models.Account{
			Email: user.Email,
		},
		Name: user.Name,
	})
	if err != nil {
		return "", err
	}

	sessionToken, err := e.jwtUtil.Sign(map[string]string{
		"email": user.Email,
		"name":  user.Name,
	}, jwt.SessionToken, time.Hour*24*30)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (e *EmailLoginService) VerifyOtp(otp entities.OtpRequest) (string, error) {
	return "", nil
}
