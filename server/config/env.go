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
		GoEnv:                    GoEnv(getEnv("GO_ENV")),
		CdnAddress:               getEnv("CDN_ADDRESS"),
		EventHubAddress:          getEnv("EVENTHUB_ADDRESS"),
		YouTubeDownloaderAddress: getEnv("YTDL_ADDRESS"),
		DankLyricsAddress:        getEnv("DANKLYRICS_ADDRESS"),
		Hostname:                 getEnv("HOST_NAME"),
		JwtSecret:                getEnv("JWT_SECRET"),
		BlobsDir:                 getEnv("BLOBS_DIR"),
		GeniusToken:              getEnv("GENIUS_CLIENT_TOKEN"),
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
		Cache: struct {
			Host     string
			Password string
		}{
			Host:     getEnv("CACHE_HOST"),
			Password: getEnv("CACHE_PASSWORD"),
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

type GoEnv string

const (
	GoEnvProd GoEnv = "prod"
	GoEnvBeta GoEnv = "beta"
	GoEnvDev  GoEnv = "dev"
	GoEnvTest GoEnv = "test"
)

type config struct {
	Port                     string
	CdnPort                  string
	EventHubPort             string
	WebPort                  string
	GoEnv                    GoEnv
	CdnAddress               string
	EventHubAddress          string
	YouTubeDownloaderAddress string
	DankLyricsAddress        string
	Hostname                 string
	JwtSecret                string
	BlobsDir                 string
	GeniusToken              string
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
	Cache struct {
		Host     string
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

func CdnAddress() string {
	cdnAddress := "https://cdn.dankmuzikk.com"
	if Env().GoEnv == "dev" {
		cdnAddress = Env().CdnAddress
	}

	return cdnAddress
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln("The \"" + key + "\" variable is missing.")
	}
	return value
}
