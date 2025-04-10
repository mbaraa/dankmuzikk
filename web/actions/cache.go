package actions

type Cache interface {
	SetRedirectPath(clientHash, path string) error
	GetRedirectPath(clientHash string) (string, error)
}
