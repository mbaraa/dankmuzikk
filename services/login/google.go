package login

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"encoding/json"
	"errors"

	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	randomState = uuid.NewString()
)

func init() {
	timer := time.NewTicker(time.Hour / 2)
	go func() {
		for range timer.C {
			randomState = uuid.NewString()
		}
	}()

}

type oauthUserInfo struct {
	Email    string `json:"email"`
	FullName string `json:"name"`
	PfpLink  string `json:"picture"`
	Locale   string `json:"locale"`
}

type GoogleLoginService struct {
	accountRepo db.CRUDRepo[models.Account]
	profileRepo db.CRUDRepo[models.Profile]
	otpRepo     db.CRUDRepo[models.EmailVerificationCode]
	jwtUtil     jwt.Manager[any]
}

func NewGoogleLoginService(
	accountRepo db.CRUDRepo[models.Account],
	profileRepo db.CRUDRepo[models.Profile],
	otpRepo db.CRUDRepo[models.EmailVerificationCode],
	jwtUtil jwt.Manager[any],
) *GoogleLoginService {
	return &GoogleLoginService{
		accountRepo: accountRepo,
		profileRepo: profileRepo,
		otpRepo:     otpRepo,
		jwtUtil:     jwtUtil,
	}
}

func (g *GoogleLoginService) Login(state, code string) (string, error) {
	googleUser, err := g.completeLoginWithGoogle(state, code)
	if err != nil {
		return "", err
	}

	account, err := g.accountRepo.GetByConds("email = ? AND is_o_auth = 1", googleUser.Email)
	if errors.Is(err, db.ErrRecordNotFound) || len(account) == 0 {
		return g.Signup(googleUser)
	}
	if err != nil {
		return "", err
	}

	profile, err := g.profileRepo.GetByConds("account_id = ?", account[0].Id)
	if err != nil {
		return "", err
	}
	profile[0].Account = account[0]
	profile[0].AccountId = account[0].Id

	verificationToken, err := g.jwtUtil.Sign(map[string]string{
		"name":  profile[0].Name,
		"email": profile[0].Account.Email,
	}, jwt.SessionToken, time.Hour*24*30)
	if err != nil {
		return "", err
	}

	return verificationToken, nil
}

func (g *GoogleLoginService) Signup(googleUser oauthUserInfo) (string, error) {
	profile := models.Profile{
		Account: models.Account{
			Email:   googleUser.Email,
			IsOAuth: true,
		},
		Name:     googleUser.FullName,
		Username: googleUser.Email,
	}
	// creating a new account will create the account underneath it.
	err := g.profileRepo.Add(&profile)
	if err != nil {
		return "", err
	}

	verificationToken, err := g.jwtUtil.Sign(map[string]string{
		"name":  profile.Name,
		"email": profile.Account.Email,
	}, jwt.SessionToken, time.Hour*24*30)
	if err != nil {
		return "", err
	}

	return verificationToken, nil
}

func (g *GoogleLoginService) completeLoginWithGoogle(state, code string) (oauthUserInfo, error) {
	if state != g.CurrentRandomState() {
		log.Errorln("[GOOGLE LOGIN]: State is invalid")
		return oauthUserInfo{}, errors.New("state is not valid")
	}

	token, err := config.GoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Exchange code is not valid")
		return oauthUserInfo{}, errors.New("Exchange code is not valid")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to fetch user info: ", err)
		return oauthUserInfo{}, err
	}
	defer response.Body.Close()

	var respBody oauthUserInfo
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to decode user info: ", err)
		return oauthUserInfo{}, err
	}

	return respBody, nil
}

func (g *GoogleLoginService) CurrentRandomState() string {
	return randomState
}
