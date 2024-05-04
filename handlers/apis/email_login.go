package apis

import (
	"context"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/services/login"
	"dankmuzikk/views/components/otp"
	"encoding/json"
	"net/http"
	"time"
)

type emailLoginApi struct {
	service login.EmailLoginService
}

func NewEmailLoginApi(service login.EmailLoginService) *emailLoginApi {
	return &emailLoginApi{service}
}

func (e *emailLoginApi) HandleEmailLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infoln(reqBody)

	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailSignup(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken, err := e.service.Signup(reqBody)
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    sessionToken,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	})
	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailOTPVerification(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.OtpRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infoln(reqBody)

	otp.VerifyOtp().Render(context.Background(), w)
}
