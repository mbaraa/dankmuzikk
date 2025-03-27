package actions

type AuthenticateUserPayload struct {
	Id       uint   `json:"-"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}

func (a *Actions) AuthenticateUser(sessionToken string) (AuthenticateUserPayload, error) {
	token, err := a.jwt.Decode(sessionToken, JwtSessionToken)
	if err != nil {
		return AuthenticateUserPayload{}, err
	}

	profile, err := a.cache.GetAuthenticatedUser(sessionToken)
	if err != nil {
		profile, err = a.app.GetProfileByAccountEmail(token.Payload.Email)
		if err != nil {
			return AuthenticateUserPayload{}, err
		}

		err = a.cache.SetAuthenticatedUser(sessionToken, profile)
		if err != nil {
			return AuthenticateUserPayload{}, err
		}
	}

	return AuthenticateUserPayload{
		Id:       profile.Id,
		Email:    profile.Account.Email,
		Name:     profile.Name,
		PfpLink:  profile.PfpLink,
		Username: profile.Username,
	}, nil
}
