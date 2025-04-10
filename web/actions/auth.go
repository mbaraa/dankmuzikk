package actions

type Profile struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}

func (a *Actions) CheckAuth(sessionToken string) error {
	return a.requests.Auth(sessionToken)
}

func (a *Actions) GetProfile(sessionToken string) (Profile, error) {
	return a.requests.GetProfile(sessionToken)
}

func (a *Actions) Logout(sessionToken string) error {
	return a.requests.Logout(sessionToken)
}

func (a *Actions) SetRedirectPath(clientHash, path string) error {
	return a.cache.SetRedirectPath(clientHash, path)
}

func (a *Actions) GetRedirectPath(clientHash string) (string, error) {
	return a.cache.GetRedirectPath(clientHash)
}
