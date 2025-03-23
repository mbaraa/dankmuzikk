package config

import (
	"dankmuzikk/log"
	"os"
)

var (
	_config = config{}
)

func initEnvVars() {
	_config = config{
		Port:                     getEnv("PORT"),
		CdnPort:                  getEnv("CDN_PORT"),
		EventHubPort:             getEnv("EVENTHUB_PORT"),
		WebPort:                  getEnv("WEB_PORT"),
		GoEnv:                    getEnv("GO_ENV"),
		CdnAddress:               getEnv("CDN_ADDRESS"),
		EventHubAddress:          getEnv("EVENTHUB_ADDRESS"),
		YouTubeDownloaderAddress: getEnv("YTDL_ADDRESS"),
		Hostname:                 getEnv("HOST_NAME"),
		JwtSecret:                getEnv("JWT_SECRET"),
		BlobsDir:                 getEnv("BLOBS_DIR"),
		Google: struct {
			ClientId     string
			ClientSecret string
		}{
			ClientId:     getEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
		},
		DB: struct {
			Name     string
			Host     string
			Username string
			Password string
		}{
			Name:     getEnv("DB_NAME"),
			Host:     getEnv("DB_HOST"),
			Username: getEnv("DB_USERNAME"),
			Password: getEnv("DB_PASSWORD"),
		},
		Smtp: struct {
			Host     string
			Port     string
			Username string
			Password string
		}{
			Host:     getEnv("SMTP_HOST"),
			Port:     getEnv("SMTP_PORT"),
			Username: getEnv("SMTP_USER"),
			Password: getEnv("SMTP_PASSWORD"),
		},
	}
}

type config struct {
	Port                     string
	CdnPort                  string
	EventHubPort             string
	WebPort                  string
	GoEnv                    string
	CdnAddress               string
	EventHubAddress          string
	YouTubeDownloaderAddress string
	Hostname                 string
	JwtSecret                string
	BlobsDir                 string
	Google                   struct {
		ClientId     string
		ClientSecret string
	}
	DB struct {
		Name     string
		Host     string
		Username string
		Password string
	}
	Smtp struct {
		Host     string
		Port     string
		Username string
		Password string
	}
}

// Env returns the thing's config values :)
func Env() config {
	return _config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln("The \"" + key + "\" variable is missing.")
	}
	return value
}
