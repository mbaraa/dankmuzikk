package entities

type LoginRequest struct {
	Email string `json:"email"`
}

type SignupRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type OtpRequest struct {
	Code string `json:"code"`
}
