package config

import (
	"dankmuzikk/log"
	"os"
)

var (
	_config = config{}
)

func init() {
	_config = config{
		Port: getEnv("PORT"),
		YouTube: struct {
			ApiServiceAccount string
			ScraperUrl        string
			MusicDir          string
		}{
			ApiServiceAccount: getEnv("GOOGLE_APPLICATION_CREDENTIALS"),
			ScraperUrl:        getEnv("YOUTUBE_SCAPER_URL"),
			MusicDir:          getEnv("YOUTUBE_MUSIC_DOWNLOAD_PATH"),
		},
	}
}

type config struct {
	Port    string
	YouTube struct {
		ApiServiceAccount string
		ScraperUrl        string
		MusicDir          string
	}
}

// Vals returns the thing's config values :)
func Vals() config {
	return _config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
