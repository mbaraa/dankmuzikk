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
	if Env().GoEnv == "prod" {
		return fmt.Sprintf("https://%s/api/login/google/callback", Env().DomainName)
	} else {
		return fmt.Sprintf("http://%s:%s/api/login/google/callback", Env().DomainName, Env().WebPort)
	}
}

func WebUrl() string {
	if Env().GoEnv == "prod" {
		return fmt.Sprintf("https://%s/", Env().DomainName)
	} else {
		return fmt.Sprintf("http://%s:%s/", Env().DomainName, Env().WebPort)
	}
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
