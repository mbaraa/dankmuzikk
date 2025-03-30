package actions

type LoginUsingEmailParams struct {
	Email string `json:"email"`
}

type LoginUsingEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) LoginUsingEmail(params LoginUsingEmailParams) (LoginUsingEmailPayload, error) {
	return a.requests.EmailLogin(params)
}

type SignupUsingEmailParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignupUsingEmailPayload struct {
	VerificationToken string `json:"verification_token"`
}

func (a *Actions) SignupUsingEmail(params SignupUsingEmailParams) (SignupUsingEmailPayload, error) {
	return a.requests.EmailSignup(params)
}

type VerifyOtpUsingEmailParams struct {
	Code  string `json:"code"`
	Token string `json:"token"`
}

type VerifyOtpUsingEmailPayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) VerifyOtpUsingEmail(params VerifyOtpUsingEmailParams) (VerifyOtpUsingEmailPayload, error) {
	return a.requests.VerifyOtp(params)
}

type LoginUsingGooglePayload struct {
	RedirectUrl string `json:"redirect_url"`
}

func (a *Actions) LoginUsingGoogle() (LoginUsingGooglePayload, error) {
	return a.requests.GoogleLogin()
}

type FinishLoginUsingGoogleParams struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type FinishLoginUsingGooglePayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) FinishLoginUsingGoogle(params FinishLoginUsingGoogleParams) (FinishLoginUsingGooglePayload, error) {
	return a.requests.GoogleFinishLogin(params)
}
