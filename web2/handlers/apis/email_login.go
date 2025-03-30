package apis

import (
	"context"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	dankerrors "dankmuzikk-web/errors"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/otp"
	"dankmuzikk-web/views/components/status"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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
	var reqBody actions.LoginUsingEmailParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := e.usecases.LoginUsingEmail(reqBody)
	if err != nil && errors.Is(err, dankerrors.ErrDifferentLoginMethodUsed) {
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
		Value:    payload.VerificationToken,
		HttpOnly: true,
		Path:     "/api/verify-otp",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour / 2),
	})
	otp.VerifyOtp().Render(context.Background(), w)
}

func (e *emailLoginApi) HandleEmailSignup(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.SignupUsingEmailParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := e.usecases.SignupUsingEmail(reqBody)
	if errors.Is(err, dankerrors.ErrAccountExists) {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		status.
			BugsBunnyError(fmt.Sprintf("An account associated with the email \"%s\" already exists", reqBody.Email)).
			Render(context.Background(), w)
		return
	}
	if err != nil {
		log.Errorf("[EMAIL LOGIN API]: Failed to sign up a new user: %+v, error: %s\n", reqBody, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.VerificationTokenKey,
		Value:    payload.VerificationToken,
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

	var reqBody actions.VerifyOtpUsingEmailParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error(err)
		status.
			BugsBunnyError("Invalid verification token").
			Render(context.Background(), w)
		return
	}

	reqBody.Token = verificationToken.Value
	payload, err := e.usecases.VerifyOtpUsingEmail(reqBody)
	if errors.Is(err, dankerrors.ErrExpiredToken) {
		status.
			BugsBunnyError("Expired verification code!").
			Render(context.Background(), w)
		return
	}
	if errors.Is(err, dankerrors.ErrInvalidVerificationCode) {
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
		Value:    payload.SessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	w.Header().Set("HX-Redirect", "/")
}
