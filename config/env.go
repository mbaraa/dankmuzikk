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
		Port:     getEnv("PORT"),
		Hostname: getEnv("HOSTNAME"),
		YouTube: struct {
			ScraperUrl string
			MusicDir   string
		}{
			ScraperUrl: getEnv("YOUTUBE_SCAPER_URL"),
			MusicDir:   getEnv("YOUTUBE_MUSIC_DOWNLOAD_PATH"),
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
	}
}

type config struct {
	Port     string
	Hostname string
	YouTube  struct {
		ScraperUrl string
		MusicDir   string
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
