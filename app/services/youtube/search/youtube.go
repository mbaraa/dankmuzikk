package search

import (
	"dankmuzikk/entities"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const (
	maxSearchResults = 7
)

// Service is an interface that represents a YouTube search behavior.
type Service interface {
	// Search searches YouTube for the given query,
	// and returns a SearchResult slice, and an occurring error.
	Search(query string) (results []entities.Song, err error)
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
		duration/(time.Hour*24), duration/time.Hour, duration/time.Minute, duration/time.Second

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
