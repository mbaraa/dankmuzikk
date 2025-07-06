package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
)

type Profile struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}

func (a *Actions) CheckAuth(sessionToken string) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/me/auth",
		Headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (a *Actions) GetProfile(ctx ActionContext) (Profile, error) {
	return requests.Do[any, Profile](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/me/profile",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
		},
	})
}

func (a *Actions) Logout(ctx ActionContext) error {
	_, err := requests.Do[any, Profile](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/me/logout",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
		},
	})
	return err
}

func (a *Actions) SetRedirectPath(clientHash, path string) error {
	return a.cache.SetRedirectPath(clientHash, path)
}

func (a *Actions) GetRedirectPath(clientHash string) (string, error) {
	return a.cache.GetRedirectPath(clientHash)
}
