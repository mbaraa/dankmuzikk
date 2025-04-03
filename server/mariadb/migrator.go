package mariadb

import (
	"dankmuzikk/app/models"
	"dankmuzikk/evy"
	"errors"
	"strings"
	"time"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Account),
		new(models.Profile),
		new(models.Song),
		new(models.Playlist),
		new(models.PlaylistSong),
		new(models.PlaylistOwner),
		new(models.History),
		new(models.PlaylistSongVoter),
		new(evy.EventPayload),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"profiles", "songs", "playlists", "event_payloads",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	return nil
}

var getDuration = durationer()

func durationer() func(strDuration string) (time.Duration, error) {
	durationSeparators := [3]rune{'s', 'm', 'h'}
	return func(strDuration string) (time.Duration, error) {
		colonsCount := 0
		for _, chr := range strDuration {
			if chr == ':' {
				colonsCount++
			}
		}
		if colonsCount > 2 {
			return 0, errors.New("invalid iso duration")
		}
		refinedStrDuration := strings.Builder{}
		for _, chr := range strDuration {
			if chr == ':' {
				refinedStrDuration.WriteRune(durationSeparators[colonsCount])
				colonsCount--
				continue
			}
			refinedStrDuration.WriteRune(chr)
		}
		refinedStrDuration.WriteRune(durationSeparators[colonsCount])

		duration, err := time.ParseDuration(refinedStrDuration.String())
		if err != nil {
			return 0, err
		}

		return duration, nil
	}
}
