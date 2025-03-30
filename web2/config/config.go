package config

import (
	"dankmuzikk-web/log"
	"os"
)

var (
	_config = config{}
)

func init() {
	_config = config{
		Port:          getEnv("WEB_PORT"),
		GoEnv:         getEnv("GO_ENV"),
		Hostname:      getEnv("HOST_NAME"),
		ServerAddress: getEnv("SERVER_ADDRESS"),
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

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
