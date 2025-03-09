package login

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
	"errors"
)

type EmailLoginService struct {
}

func NewEmailLoginService() *EmailLoginService {
	return &EmailLoginService{}
}

func (e *EmailLoginService) Login(user entities.LoginRequest) (string, error) {
	respBody, err := requests.PostRequest[entities.LoginRequest, map[string]string]("/v1/login/email", user)
	if err != nil {
		return "", err
	}

	token := respBody["verification_token"]
	if token == "" {
		return "", errors.New("oopsie")
	}

	return token, nil
}

func (e *EmailLoginService) Signup(user entities.SignupRequest) (string, error) {
	respBody, err := requests.PostRequest[entities.SignupRequest, map[string]string]("/v1/signup/email", user)
	if err != nil {
		return "", err
	}

	token := respBody["verification_token"]
	if token == "" {
		return "", errors.New("oopsie")
	}

	return token, nil
}

func (e *EmailLoginService) VerifyOtp(otp entities.OtpRequest) (string, error) {
	respBody, err := requests.PostRequest[entities.OtpRequest, map[string]string]("/v1/verify-otp", otp)
	if err != nil {
		return "", err
	}

	token := respBody["session_token"]
	if token == "" {
		return "", errors.New("oopsie")
	}

	return token, nil
}
