package login

type ErrAccountNotFound struct{}

func (e ErrAccountNotFound) Error() string {
	return "account-not-found"
}

func (e ErrAccountNotFound) ClientStatusCode() int {
	return 404
}

func (e ErrAccountNotFound) ExtraData() map[string]any {
	return nil
}

func (e ErrAccountNotFound) ExposeToClients() bool {
	return true
}

type ErrProfileNotFound struct{}

func (e ErrProfileNotFound) Error() string {
	return "profile-not-found"
}

func (e ErrProfileNotFound) ClientStatusCode() int {
	return 404
}

func (e ErrProfileNotFound) ExtraData() map[string]any {
	return nil
}

func (e ErrProfileNotFound) ExposeToClients() bool {
	return true
}

type ErrAccountExists struct{}

func (e ErrAccountExists) Error() string {
	return "account-already-exists"
}

func (e ErrAccountExists) ClientStatusCode() int {
	return 409
}

func (e ErrAccountExists) ExtraData() map[string]any {
	return nil
}

func (e ErrAccountExists) ExposeToClients() bool {
	return true
}

type ErrExpiredVerificationCode struct{}

func (e ErrExpiredVerificationCode) Error() string {
	return "verification-code-expired"
}

func (e ErrExpiredVerificationCode) ClientStatusCode() int {
	return 400
}

func (e ErrExpiredVerificationCode) ExtraData() map[string]any {
	return nil
}

func (e ErrExpiredVerificationCode) ExposeToClients() bool {
	return true
}

type ErrInvalidVerificationCode struct{}

func (e ErrInvalidVerificationCode) Error() string {
	return "invalid-verification-code"
}

func (e ErrInvalidVerificationCode) ClientStatusCode() int {
	return 400
}

func (e ErrInvalidVerificationCode) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidVerificationCode) ExposeToClients() bool {
	return true
}

type ErrDifferentLoginMethod struct{}

func (e ErrDifferentLoginMethod) Error() string {
	return "different-login-method-used"
}

func (e ErrDifferentLoginMethod) ClientStatusCode() int {
	return 409
}

func (e ErrDifferentLoginMethod) ExtraData() map[string]any {
	return nil
}

func (e ErrDifferentLoginMethod) ExposeToClients() bool {
	return true
}
