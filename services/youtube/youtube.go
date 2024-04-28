package youtube

import (
	"dankmuzikk/log"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	maxSearchResults = 7
)

func init() {
	// youtube api
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalln(log.ErrorLevel, "[YOUTUBE SERVICE] Missing Google API Service Account File")
	}
	// youtube scraper 👀
	if os.Getenv("YOUTUBE_SCAPER_URL") == "" {
		log.Fatalln(log.ErrorLevel, "[YOUTUBE SERVICE] Missing YouTube Scrapper URL")
	}
}

type SearchResult struct {
	Id           string
	Title        string
	ChannelTitle string // or artist idk
	Description  string
	ThumbnailUrl string
	Duration     string
}

// YouTubeSearcher is an interface that represents a youtube search behavior.
type YouTubeSearcher interface {
	// Search searches youtube for the given query,
	// and returns a SearchResult slice, and an occurring error.
	Search(query string) (results []SearchResult, err error)
}

func getTime(isoDuration string) (string, error) {
	startIdx := 0
	for i, chr := range isoDuration {
		if unicode.IsDigit(chr) {
			startIdx = i
			break
		}
	}
	duration, err := time.ParseDuration(strings.ToLower(isoDuration[startIdx:]))
	if err != nil {
		return "", err
	}
	days, hours, mins, secs :=
		duration/(time.Hour*24), (duration / time.Hour), duration/time.Minute, duration/time.Second

	builder := strings.Builder{}
	if days > 0 {
		builder.WriteString(formatNumber(int(days)) + ":")
	}
	if hours > 0 {
		builder.WriteString(formatNumber(int(hours%60)) + ":")
	}
	if mins > 0 {
		builder.WriteString(formatNumber(int(mins%60)) + ":")
	}
	if secs > 0 {
		if days == 0 && hours == 0 && mins == 0 {
			builder.WriteString(formatNumber(int(secs%60)) + "s")
		} else {
			builder.WriteString(formatNumber(int(secs % 60)))
		}
	}

	return builder.String(), nil
}

func formatNumber(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprint(n)
}
