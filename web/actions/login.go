package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
)

type LoginUsingEmailParams struct {
	Email string `json:"email"`
}

type LoginUsingEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) LoginUsingEmail(params LoginUsingEmailParams) (LoginUsingEmailPayload, error) {
	return requests.Do[LoginUsingEmailParams, LoginUsingEmailPayload](requests.Config[LoginUsingEmailParams]{
		Method:   http.MethodPost,
		Endpoint: "/v1/login/email",
		Body:     params,
	})
}

type SignupUsingEmailParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignupUsingEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) SignupUsingEmail(params SignupUsingEmailParams) (SignupUsingEmailPayload, error) {
	return requests.Do[SignupUsingEmailParams, SignupUsingEmailPayload](requests.Config[SignupUsingEmailParams]{
		Method:   http.MethodPost,
		Endpoint: "/v1/signup/email",
		Body:     params,
	})
}

type VerifyOtpUsingEmailParams struct {
	Code  string `json:"code"`
	Token string `json:"token"`
}

type VerifyOtpUsingEmailPayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) VerifyOtpUsingEmail(params VerifyOtpUsingEmailParams) (VerifyOtpUsingEmailPayload, error) {
	return requests.Do[VerifyOtpUsingEmailParams, VerifyOtpUsingEmailPayload](requests.Config[VerifyOtpUsingEmailParams]{
		Method:   http.MethodPost,
		Endpoint: "/v1/verify-otp",
		Body:     params,
	})
}

type LoginUsingGooglePayload struct {
	RedirectUrl string `json:"redirect_url"`
}

func (a *Actions) LoginUsingGoogle() (LoginUsingGooglePayload, error) {
	return requests.Do[any, LoginUsingGooglePayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/login/google",
	})
}

type FinishLoginUsingGoogleParams struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type FinishLoginUsingGooglePayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) FinishLoginUsingGoogle(params FinishLoginUsingGoogleParams) (FinishLoginUsingGooglePayload, error) {
	return requests.Do[FinishLoginUsingGoogleParams, FinishLoginUsingGooglePayload](requests.Config[FinishLoginUsingGoogleParams]{
		Method:   http.MethodPost,
		Endpoint: "/v1/login/google/callback",
		Body:     params,
	})
}
