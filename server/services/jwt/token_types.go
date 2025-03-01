package jwt

// Subject represents the token's subject
type Subject = string

const (
	// SessionToken used to verify that the user is logged in correctly and can access the good stuff.
	SessionToken Subject = "SESSION_TOKEN"
	// VerificationToken used to verify verification code (on sign-in or sign-up).
	VerificationToken Subject = "OTP_TOKEN"
)
