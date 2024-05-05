package apis

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/login"
	"dankmuzikk/views/components/otp"
	"encoding/json"
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
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     handlers.VerificationTokenKey,
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
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     handlers.VerificationTokenKey,
		Value:    verificationToken,
		HttpOnly: true,
		Path:     "/api/verify-otp",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour / 2),
	})
	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailOTPVerification(w http.ResponseWriter, r *http.Request) {
	verificationToken, err := r.Cookie(handlers.VerificationTokenKey)
	if err != nil {
		w.Write([]byte("Invalid verification token"))
		return
	}
	if verificationToken.Expires.After(time.Now().UTC()) {
		w.Write([]byte("Expired verification token"))
		return
	}

	var reqBody entities.OtpRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken, err := e.service.VerifyOtp(verificationToken.Value, reqBody)
	// TODO: specify errors further suka
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     handlers.SessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	})

	w.Header().Set("HX-Redirect", "/")
}
