package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
)

type emailLoginApi struct {
	usecases *actions.Actions
}

func NewEmailLoginApi(usecases *actions.Actions) *emailLoginApi {
	return &emailLoginApi{
		usecases: usecases,
	}
}

func (e *emailLoginApi) HandleEmailLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.LoginWithEmailParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.LoginWithEmail(reqBody)
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *emailLoginApi) HandleEmailSignup(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.SignupWithEmailParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.SignupWithEmail(reqBody)
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *emailLoginApi) HandleEmailOTPVerification(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.VerifyOtpParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.VerifyOtp(reqBody)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
