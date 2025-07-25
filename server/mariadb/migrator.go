package mariadb

import (
	"dankmuzikk/app/models"
	"dankmuzikk/evy"
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
		new(models.FavoriteSong),
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
