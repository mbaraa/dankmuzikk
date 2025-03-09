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
		Port:      getEnv("PORT"),
		GoEnv:     getEnv("GO_ENV"),
		Hostname:  getEnv("HOSTNAME"),
		JwtSecret: getEnv("JWT_SECRET"),
		YouTube: struct {
			DownloaderUrl string
			MuzikkDir     string
		}{
			DownloaderUrl: getEnv("YOUTUBE_DOWNLOADER_URL"),
			MuzikkDir:     getEnv("YOUTUBE_MUSIC_DOWNLOAD_PATH"),
		},
		Google: struct {
			ApiServiceAccount string
			ClientId          string
			ClientSecret      string
		}{
			ApiServiceAccount: getEnv("GOOGLE_APPLICATION_CREDENTIALS"),
			ClientId:          getEnv("GOOGLE_CLIENT_ID"),
			ClientSecret:      getEnv("GOOGLE_CLIENT_SECRET"),
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
	Port      string
	GoEnv     string
	Hostname  string
	JwtSecret string
	YouTube   struct {
		DownloaderUrl string
		MuzikkDir     string
	}
	Google struct {
		ApiServiceAccount string
		ClientId          string
		ClientSecret      string
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
