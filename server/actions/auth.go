package actions

import (
	"context"
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/log"
	"dankmuzikk/nanoid"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	verificationTokenTtlMinutes = 30
	sessionTokenTtlDays         = 60
)

func (a *Actions) AuthenticateAccount(sessionToken string) (models.Account, error) {
	token, err := a.jwt.Decode(sessionToken, JwtSessionToken)
	if err != nil {
		return models.Account{}, err
	}

	account, err := a.cache.GetAuthenticatedAccount(sessionToken)
	if err != nil {
		account, err = a.app.GetAccountByEmail(token.Payload.Email)
		if err != nil {
			return models.Account{}, err
		}

		err = a.cache.SetAuthenticatedAccount(sessionToken, account)
		if err != nil {
			return models.Account{}, err
		}
	}

	return account, nil
}

type GetProfilePayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}

func (a *Actions) GetProfile(email string) (GetProfilePayload, error) {
	profile, err := a.app.GetProfileByAccountEmail(email)
	if err != nil {
		return GetProfilePayload{}, err
	}

	return GetProfilePayload{
		Email:    profile.Account.Email,
		Name:     profile.Name,
		PfpLink:  profile.PfpLink,
		Username: profile.Username,
	}, nil
}

type TokenPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (t TokenPayload) Valid() bool {
	return t.Name != "" && t.Email != "" && t.Username != ""
}

type LoginWithEmailParams struct {
	Email string `json:"email"`
}

type LoginWithEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) LoginWithEmail(params LoginWithEmailParams) (LoginWithEmailPayload, error) {
	profile, err := a.app.GetProfileByAccountEmail(params.Email)
	if err != nil {
		return LoginWithEmailPayload{}, err
	}
	if profile.Account.IsOAuth {
		return LoginWithEmailPayload{}, &app.ErrDifferentLoginMethod{}
	}

	verificationToken, err := a.jwt.Sign(TokenPayload{
		Name:     profile.Name,
		Email:    profile.Account.Email,
		Username: profile.Account.Email,
	}, JwtVerificationToken, time.Minute*verificationTokenTtlMinutes)
	if err != nil {
		return LoginWithEmailPayload{}, err
	}

	err = a.sendOtp(profile)
	if err != nil {
		return LoginWithEmailPayload{}, err
	}

	return LoginWithEmailPayload{
		VerificationToken: verificationToken,
	}, nil
}

type SignupWithEmailParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignupWithEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) SignupWithEmail(params SignupWithEmailParams) (SignupWithEmailPayload, error) {
	profile, err := a.app.CreateNoOAuthUser(app.CreateNoOAuthUserArgs{
		Email: params.Email,
		Name:  params.Name,
	})
	if err != nil {
		return SignupWithEmailPayload{}, err
	}

	verificationToken, err := a.jwt.Sign(TokenPayload{
		Name:     params.Name,
		Email:    params.Email,
		Username: params.Email,
	}, JwtVerificationToken, time.Minute*verificationTokenTtlMinutes)
	if err != nil {
		return SignupWithEmailPayload{}, err
	}

	err = a.sendOtp(profile)
	if err != nil {
		return SignupWithEmailPayload{}, err
	}

	return SignupWithEmailPayload{
		VerificationToken: verificationToken,
	}, nil
}

type VerifyOtpParams struct {
	Code  string `json:"code"`
	Token string `json:"token"`
}

type VerifyOtpPayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) VerifyOtp(params VerifyOtpParams) (VerifyOtpPayload, error) {
	tokeeen, err := a.jwt.Decode(params.Token, JwtVerificationToken)
	if err != nil {
		return VerifyOtpPayload{}, err
	}

	if !tokeeen.Payload.Valid() {
		return VerifyOtpPayload{}, &app.ErrInvalidVerificationToken{}
	}

	account, err := a.app.GetAccountByEmail(tokeeen.Payload.Email)
	if err != nil {
		return VerifyOtpPayload{}, err
	}

	otp, err := a.app.GetOtpForAccount(account.Id)
	if err != nil {
		return VerifyOtpPayload{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(otp), []byte(params.Code))
	if err != nil {
		return VerifyOtpPayload{}, &app.ErrInvalidVerificationToken{}
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Name:     tokeeen.Payload.Name,
		Email:    tokeeen.Payload.Email,
		Username: tokeeen.Payload.Email,
	}, JwtSessionToken, time.Hour*24*sessionTokenTtlDays)
	if err != nil {
		return VerifyOtpPayload{}, err
	}

	return VerifyOtpPayload{
		SessionToken: sessionToken,
	}, nil
}

func (a *Actions) sendOtp(profile models.Profile) error {
	otp := generateOtp()

	hashedOtp, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = a.app.CreateOtp(profile.AccountId, string(hashedOtp))
	if err != nil {
		return err
	}

	if config.Env().GoEnv == "dev" {
		log.Warningf("\n\n################################\n#### OTP SUKA: %s\n################################\n\n", otp)
	} else {
		err = a.mailer.SendOtpEmail(profile, otp)
		if err != nil {
			return err
		}
	}

	return nil
}

type LoginWithGoogleParams struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type LoginWithGooglePayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) LoginWithGoogle(params LoginWithGoogleParams) (LoginWithGooglePayload, error) {
	googleUser, err := a.completeLoginWithGoogle(params)
	if err != nil {
		return LoginWithGooglePayload{}, err
	}

	profile, err := a.app.GetProfileByAccountEmail(googleUser.Email)
	if _, ok := err.(*app.ErrNotFound); ok {
		return a.signupWithGoogle(googleUser)
	}
	if !profile.Account.IsOAuth {
		return LoginWithGooglePayload{}, &app.ErrDifferentLoginMethod{}
	}
	if err != nil {
		return LoginWithGooglePayload{}, err
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Name:     profile.Name,
		Email:    profile.Account.Email,
		Username: profile.Account.Email,
	}, JwtSessionToken, time.Hour*24*sessionTokenTtlDays)
	if err != nil {
		return LoginWithGooglePayload{}, err
	}

	return LoginWithGooglePayload{
		SessionToken: sessionToken,
	}, nil
}

func (a *Actions) signupWithGoogle(params googleOAuthUserInfo) (LoginWithGooglePayload, error) {
	profile, err := a.app.CreateOAuthUser(app.CreateOAuthUserArgs{
		Email:   params.Email,
		Name:    params.FullName,
		PfpLink: params.PfpLink,
	})
	if err != nil {
		return LoginWithGooglePayload{}, err
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Name:     profile.Name,
		Email:    profile.Account.Email,
		Username: profile.Account.Email,
	}, JwtSessionToken, time.Hour*24*sessionTokenTtlDays)
	if err != nil {
		return LoginWithGooglePayload{}, err
	}

	return LoginWithGooglePayload{
		SessionToken: sessionToken,
	}, nil
}

type googleOAuthUserInfo struct {
	Email    string `json:"email"`
	FullName string `json:"name"`
	PfpLink  string `json:"picture"`
	Locale   string `json:"locale"`
}

func (a *Actions) completeLoginWithGoogle(params LoginWithGoogleParams) (googleOAuthUserInfo, error) {
	if params.State != a.CurrentGoogleLoginRandomState() {
		log.Errorln("[GOOGLE LOGIN]: State is invalid")
		return googleOAuthUserInfo{}, errors.New("state is not valid")
	}

	token, err := config.GoogleOAuthConfig().Exchange(context.Background(), params.Code)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Exchange code is not valid")
		return googleOAuthUserInfo{}, errors.New("Exchange code is not valid")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to fetch user info: ", err)
		return googleOAuthUserInfo{}, err
	}
	defer response.Body.Close()

	var respBody googleOAuthUserInfo
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to decode user info: ", err)
		return googleOAuthUserInfo{}, err
	}

	return respBody, nil
}

func (a *Actions) CurrentGoogleLoginRandomState() string {
	state, err := a.cache.GetGoogleLoginState()
	if err != nil {
		state = nanoid.NewWithLength(32)
		err = a.cache.SetGoogleLoginState(state)
		if err != nil {
			return state
		}
	}

	return state
}

func (a *Actions) InvalidateAuthenticatedAccount(token string) error {
	return a.cache.InvalidateAuthenticatedAccount(token)
}

func generateOtp() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	n := r.Intn(1_000_000_000-100001) + 100001
	return fmt.Sprint(n)[:6]
}
