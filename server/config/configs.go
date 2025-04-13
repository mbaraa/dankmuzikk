package config

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

// GoogleLoginRedirectUrl redirects to frontend after logging with suka Google.
func GoogleLoginRedirectUrl() string {
	return fmt.Sprintf("%s/api/login/google/callback", Env().Hostname)
}

func initGoogleConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  GoogleLoginRedirectUrl(),
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
