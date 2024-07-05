package apis

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/entities"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"dankmuzikk/services/login"
	"dankmuzikk/views/components/otp"
	"dankmuzikk/views/components/status"
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
		// w.WriteHeader(http.StatusInternalServerError)
		status.
			BugsBunnyError("This account uses Google Auth to login!").
			Render(context.Background(), w)
		return
	} else if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
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
		// w.WriteHeader(http.StatusConflict)
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
		// w.Write([]byte("Invalid verification token"))
		status.
			BugsBunnyError("Invalid verification token").
			Render(context.Background(), w)
		return
	}
	if verificationToken.Expires.After(time.Now().UTC()) {
		// w.Write([]byte("Expired verification token"))
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

	sessionToken, err := e.service.VerifyOtp(verificationToken.Value, reqBody)
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
		// w.WriteHeader(http.StatusInternalServerError)
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
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	})

	w.Header().Set("HX-Redirect", "/")
}
