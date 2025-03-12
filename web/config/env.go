package config

import (
	"dankmuzikk-web/log"
	"fmt"
	"os"
)

var (
	_config = config{}
)

func initEnvVars() {
	_config = config{
		Port:          getEnv("WEB_PORT"),
		GoEnv:         getEnv("GO_ENV"),
		Hostname:      getEnv("HOST_NAME"),
		ServerAddress: getEnv("BACKEND_ADDRESS"),
		CdnAddress:    getEnv("CDN_ADDRESS"),
	}
}

type config struct {
	Port          string
	GoEnv         string
	Hostname      string
	ServerAddress string
	CdnAddress    string
}

// Env returns the thing's config values :)
func Env() config {
	return _config
}

// GetRequestUrl returns a full path of the request with the backend's URL.
// TODO: make a unified requests package or something.
func GetRequestUrl(path string) string {
	return fmt.Sprintf("%s%s", Env().ServerAddress, path)
}

func GetCdnUrl() string {
	return fmt.Sprintf("%s", Env().CdnAddress)
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
