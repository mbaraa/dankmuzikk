package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func initGoogleConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  Env().Hostname + "/api/login/google/callback",
		ClientID:     Env().Google.ClientId,
		ClientSecret: Env().Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleOAuthConfig() *oauth2.Config {
	return googleOauthConfig
}
