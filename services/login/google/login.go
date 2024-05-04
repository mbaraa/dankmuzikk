package google

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/log"
	"dankmuzikk/services/jwt"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	randomState = uuid.NewString()
	jwtUtil     = jwt.NewJWTImpl()
)

func init() {
	timer := time.NewTicker(time.Hour / 2)
	go func() {
		for range timer.C {
			randomState = uuid.NewString()
		}
	}()

}

func CurrentRandomState() string {
	return randomState
}

type oauthUserInfo struct {
	Email    string `json:"email"`
	FullName string `json:"name"`
	PfpLink  string `json:"picture"`
	Locale   string `json:"locale"`
}

func CompleteLoginWithGoogle(state, code string) (string, error) {
	if state != CurrentRandomState() {
		log.Errorln("[GOOGLE LOGIN]: State is invalid")
		return "", errors.New("state is not valid")
	}

	token, err := config.GoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Exchange code is not valid")
		return "", errors.New("Exchange code is not valid")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to fetch user info: ", err)
		return "", err
	}
	defer response.Body.Close()

	var respBody oauthUserInfo
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to decode user info: ", err)
		return "", err
	}

	sessionToken, err := jwtUtil.Sign(respBody, jwt.SessionToken, time.Hour*24*30)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to generate jwt: ", err)
		return "", err
	}

	return sessionToken, nil
}
