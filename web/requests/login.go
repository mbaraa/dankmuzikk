package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) EmailLogin(params actions.LoginUsingEmailParams) (actions.LoginUsingEmailPayload, error) {
	return makeRequest[actions.LoginUsingEmailParams, actions.LoginUsingEmailPayload](makeRequestConfig[actions.LoginUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/login/email",
		body:     params,
	})
}

func (r *Requests) EmailSignup(params actions.SignupUsingEmailParams) (actions.SignupUsingEmailPayload, error) {
	return makeRequest[actions.SignupUsingEmailParams, actions.SignupUsingEmailPayload](makeRequestConfig[actions.SignupUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/signup/email",
		body:     params,
	})
}

func (r *Requests) VerifyOtp(params actions.VerifyOtpUsingEmailParams) (actions.VerifyOtpUsingEmailPayload, error) {
	return makeRequest[actions.VerifyOtpUsingEmailParams, actions.VerifyOtpUsingEmailPayload](makeRequestConfig[actions.VerifyOtpUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/verify-otp",
		body:     params,
	})
}

func (r *Requests) GoogleLogin() (actions.LoginUsingGooglePayload, error) {
	return makeRequest[any, actions.LoginUsingGooglePayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/login/google",
	})
}

func (r *Requests) GoogleFinishLogin(params actions.FinishLoginUsingGoogleParams) (actions.FinishLoginUsingGooglePayload, error) {
	return makeRequest[actions.FinishLoginUsingGoogleParams, actions.FinishLoginUsingGooglePayload](makeRequestConfig[actions.FinishLoginUsingGoogleParams]{
		method:   http.MethodPost,
		endpoint: "/v1/login/google/callback",
		body:     params,
	})
}
