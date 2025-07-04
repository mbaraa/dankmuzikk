package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) Auth(sessionToken string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/auth",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (r *Requests) GetProfile(sessionToken string) (actions.Profile, error) {
	return makeRequest[any, actions.Profile](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/profile",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
}

func (r *Requests) Logout(sessionToken string) error {
	_, err := makeRequest[any, actions.Profile](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/logout",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}
