package apis

import (
	"context"
	"dankmuzikk-web/config"
	"dankmuzikk-web/entities"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/login"
	"dankmuzikk-web/views/components/otp"
	"dankmuzikk-web/views/components/status"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type emailLoginApi struct {
	service *login.EmailLoginService
}

func NewEmailLoginApi(service *login.EmailLoginService) *emailLoginApi {
	return &emailLoginApi{service}
}

func (e *emailLoginApi) HandleEmailLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verificationToken, err := e.service.Login(reqBody)
	if err != nil && errors.Is(err, login.ErrDifferentLoginMethod) {
		log.Errorf("[EMAIL LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		status.
			BugsBunnyError("This account uses Google Auth to login!").
			Render(context.Background(), w)
		return
	} else if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		status.
			BugsBunnyError(fmt.Sprintf("No account associated with the email \"%s\" was found", reqBody.Email)).
			Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.VerificationTokenKey,
		Value:    verificationToken,
		HttpOnly: true,
		Path:     "/api/verify-otp",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour / 2),
	})
	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailSignup(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verificationToken, err := e.service.Signup(reqBody)
	if errors.Is(err, login.ErrAccountExists) {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		status.
			BugsBunnyError(fmt.Sprintf("An account associated with the email \"%s\" already exists", reqBody.Email)).
			Render(context.Background(), w)
		return
	}
	if errors.Is(err, login.ErrAccountNotFound) || errors.Is(err, login.ErrProfileNotFound) {

	}
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.VerificationTokenKey,
		Value:    verificationToken,
		HttpOnly: true,
		Path:     "/api/verify-otp",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour / 2),
	})
	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailOTPVerification(w http.ResponseWriter, r *http.Request) {
	verificationToken, err := r.Cookie(auth.VerificationTokenKey)
	if err != nil {
		status.
			BugsBunnyError("Invalid verification token").
			Render(context.Background(), w)
		return
	}
	if verificationToken.Expires.After(time.Now().UTC()) {
		status.
			BugsBunnyError("Expired verification token").
			Render(context.Background(), w)
		return
	}

	var reqBody entities.OtpRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error(err)
		// w.WriteHeader(http.StatusBadRequest)
		status.
			BugsBunnyError("Invalid verification token").
			Render(context.Background(), w)
		return
	}

	reqBody.Token = verificationToken.Value
	sessionToken, err := e.service.VerifyOtp(reqBody)
	if errors.Is(err, login.ErrExpiredVerificationCode) {
		status.
			BugsBunnyError("Expired verification code!").
			Render(context.Background(), w)
		return
	}
	if errors.Is(err, login.ErrInvalidVerificationCode) {
		status.
			BugsBunnyError("Invalid verification code!").
			Render(context.Background(), w)
		return
	}
	if err != nil {
		log.Error(err)
		status.
			BugsBunnyError("Something went wrong...").
			Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	w.Header().Set("HX-Redirect", "/")
}
