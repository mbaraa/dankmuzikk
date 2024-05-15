package migrator

import (
	"dankmuzikk/db"
	"dankmuzikk/models"
)

func Migrate() error {
	dbConn, err := db.Connector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Account),
		new(models.Profile),
		new(models.EmailVerificationCode),
		new(models.Song),
		new(models.Playlist),
		new(models.PlaylistSong),
		new(models.PlaylistOwner),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"profiles", "songs", "playlists",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci").Error
		if err != nil {
			return err
		}
	}

	return nil
}
