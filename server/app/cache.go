package app

type Cache interface {
	CreateOtp(accountId uint, otp string) error
	GetOtpForAccount(id uint) (string, error)
}
