package google

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/log"
	"encoding/json"
	"errors"
	"fmt"
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

func CurrentRandomState() string {
	return randomState
}

type oauthUserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	FullName  string `json:"name"`
	PfpLink   string `json:"picture"`
	Locale    string `json:"locale"`
}

func CompleteLoginWithGoogle(state, code string) error {
	if state != CurrentRandomState() {
		log.Errorln("[GOOGLE LOGIN]: State is invalid")
		return errors.New("state is not valid")
	}

	token, err := config.GoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Exchange code is not valid")
		return errors.New("Exchange code is not valid")
	}

	// Use the token to get user information from Google
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to fetch user info: ", err)
		return err
	}
	defer response.Body.Close()

	// Decode JSON response
	var respBody oauthUserInfo
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN]: Failed to decode user info: ", err)
		return err
	}

	fmt.Printf("UserInfo: %+v\n", respBody)

	return nil
}
