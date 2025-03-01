package apis

import (
	"dankmuzikk/app/entities"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"dankmuzikk/services/login"
	"encoding/json"
	"net/http"
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
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": verificationToken,
	})
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
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": verificationToken,
	})
}

func (e *emailLoginApi) HandleEmailOTPVerification(w http.ResponseWriter, r *http.Request) {
	verificationToken, err := r.Cookie(auth.VerificationTokenKey)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	var reqBody entities.OtpRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	sessionToken, err := e.service.VerifyOtp(verificationToken.Value, reqBody)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": sessionToken,
	})
}
