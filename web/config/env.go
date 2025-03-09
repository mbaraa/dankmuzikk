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
		Port:        getEnv("WEB_PORT"),
		BackendPort: getEnv("PORT"),
		CdnPort:     getEnv("CDN_PORT"),
		GoEnv:       getEnv("GO_ENV"),
		DomainName:  getEnv("DOMAIN_NAME"),
		JwtSecret:   getEnv("JWT_SECRET"),
	}
}

type config struct {
	Port        string
	BackendPort string
	CdnPort     string
	GoEnv       string
	DomainName  string
	JwtSecret   string
}

// Env returns the thing's config values :)
func Env() config {
	return _config
}

// GetRequestUrl returns a full path of the request with the backend's URL.
// TODO: make a unified requests package or something.
func GetRequestUrl(path string) string {
	if Env().GoEnv == "prod" {
		return fmt.Sprintf("https://api.%s%s", Env().DomainName, path)
	} else {
		return fmt.Sprintf("http://%s:%s%s", Env().DomainName, Env().BackendPort, path)
	}
}

func GetCdnUrl() string {
	if Env().GoEnv == "prod" {
		return fmt.Sprintf("https://cdn.%s/muzikkx", Env().DomainName)
	} else if Env().GoEnv == "dev" {
		return fmt.Sprintf("http://%s:20251/muzikkx", Env().DomainName)
	} else {
		return fmt.Sprintf("http://%s:%s/muzikkx", Env().DomainName, Env().CdnPort)
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
